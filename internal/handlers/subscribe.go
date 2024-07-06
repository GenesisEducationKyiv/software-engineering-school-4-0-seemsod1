package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/seemsod1/api-project/pkg/forms"
	"github.com/seemsod1/api-project/pkg/timezone"

	"github.com/seemsod1/api-project/internal/storage/dbrepo"

	"go.uber.org/zap"

	"github.com/go-chi/render"
	"github.com/seemsod1/api-project/internal/models"
)

type Subscriber interface {
	AddSubscriber(subscriber models.Subscriber) error
	RemoveSubscriber(email string) error
	GetSubscribersWithTimezone(timezone int) ([]string, error)
	GetSubscribers() ([]string, error)
}

// Subscribe subscribes a user to the newsletter
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	email, err := parseEmail(r)
	if err != nil {
		m.Logger.Error("Invalid email", zap.Error(err))
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}
	offset, err := timezone.ProcessTimezoneHeader(r)
	if err != nil {
		m.Logger.Error("Invalid timezone", zap.Error(err))
		http.Error(w, "Invalid timezone", http.StatusBadRequest)
		return
	}

	if err = m.Subscriber.AddSubscriber(models.Subscriber{
		Email:    email,
		Timezone: offset,
	}); err != nil {
		if errors.Is(err, dbrepo.ErrorDuplicateSubscription) {
			m.Logger.Errorf("%s: %s", email, "Already exists")
			http.Error(w, "Already exists", http.StatusConflict)
			return
		}
		m.Logger.Errorf("adding subscriber: %w", err)
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Successfully subscribed"})
}

// Unsubscribe unsubscribes a user from the newsletter
func (m *Repository) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	email, err := parseEmail(r)
	if err != nil {
		m.Logger.Error("Invalid email", zap.Error(err))
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	if err = m.Subscriber.RemoveSubscriber(email); err != nil {
		if errors.Is(err, dbrepo.ErrorNonExistentSubscription) {
			m.Logger.Errorf("%s: %s", email, "Not found")
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		m.Logger.Errorf("removing subscriber: %w", err)
		http.Error(w, "Failed to unsubscribe", http.StatusInternalServerError)
		return
	}

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
