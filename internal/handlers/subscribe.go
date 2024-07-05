package handlers

import (
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/render"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/forms"
	"github.com/seemsod1/api-project/internal/models"
	"github.com/seemsod1/api-project/internal/timezone"
)

type Subscriber interface {
	AddSubscriber(subscriber models.Subscriber) error
	GetSubscribersWithTimezone(timezone int) ([]string, error)
	GetSubscribers() ([]string, error)
}

// Subscribe subscribes a user to the newsletter
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		m.Logger.Errorf("parsing form: %w", err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("email")
	form.IsEmail("email")

	if !form.Valid() {
		m.Logger.Error("Invalid email", zap.String("email", email))
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
		if errors.Is(err, customerrors.ErrDuplicatedKey) {
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
