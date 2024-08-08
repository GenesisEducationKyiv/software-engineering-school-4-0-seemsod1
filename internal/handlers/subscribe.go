package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"

	customerrepo "github.com/seemsod1/api-project/internal/customer/repository"
	subscriberrepo "github.com/seemsod1/api-project/internal/subscriber/repository"

	"github.com/seemsod1/api-project/pkg/forms"
	"github.com/seemsod1/api-project/pkg/timezone"

	"go.uber.org/zap"

	"github.com/go-chi/render"
)

var (
	unsubscribeSuccessTotal      = metrics.NewCounter("unsubscribe_success_total")
	subscribeStatusConflictTotal = metrics.NewCounter("subscribe_status_conflict_total")
	subscribeBadRequestTotal     = metrics.NewCounter("subscribe_bad_request_total")
)

type customer interface {
	StartTransaction(email string, timezone int) error
}

type subscriber interface {
	RemoveSubscriber(email string) error
}

// Subscribe subscribes a user to the newsletter
func (m *Handlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	email, err := parseEmail(r)
	if err != nil {
		m.Logger.Error("Invalid email", zap.Error(err))
		subscribeBadRequestTotal.Inc()
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	offset, err := timezone.ProcessTimezoneHeader(r)
	if err != nil {
		m.Logger.Error("Invalid timezone", zap.Error(err))
		subscribeBadRequestTotal.Inc()
		http.Error(w, "Invalid timezone", http.StatusBadRequest)
		return
	}

	if er := m.Customer.StartTransaction(email, offset); er != nil {
		if errors.Is(er, customerrepo.ErrorDuplicateCustomer) {
			m.Logger.Error("Already exists", zap.String("email", email), zap.Error(er))
			subscribeStatusConflictTotal.Inc()
			http.Error(w, "Already exists", http.StatusConflict)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Thank you for subscribing!"})
}

// Unsubscribe unsubscribes a user from the newsletter
func (m *Handlers) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	email, err := parseEmail(r)
	if err != nil {
		m.Logger.Error("Invalid email", zap.Error(err))
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	if err = m.Subscriber.RemoveSubscriber(email); err != nil {
		if errors.Is(err, subscriberrepo.ErrorNonExistentSubscription) {
			m.Logger.Error("Not found", zap.String("email", email), zap.Error(err))
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		m.Logger.Error("removing subscriber", zap.Error(err))
		http.Error(w, "Failed to unsubscribe", http.StatusInternalServerError)
		return
	}

	unsubscribeSuccessTotal.Inc()
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Successfully unsubscribed"})
}

func parseEmail(r *http.Request) (string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return "", fmt.Errorf("unable to parse form")
	}

	email := r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("email")
	form.IsEmail("email")

	if !form.Valid() {
		return "", fmt.Errorf("invalid email")
	}

	return email, nil
}
