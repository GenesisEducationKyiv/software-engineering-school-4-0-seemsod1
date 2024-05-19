package handlers

import (
	"errors"
	"github.com/go-chi/render"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/forms"
	"github.com/seemsod1/api-project/internal/helpers"
	"github.com/seemsod1/api-project/internal/models"
	"net/http"
)

// Subscribe subscribes a user to the newsletter
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("email")
	form.IsEmail("email")

	if !form.Valid() {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	offset, err := helpers.ProcessTimezoneHeader(r)
	if err != nil {
		http.Error(w, "Invalid timezone", http.StatusBadRequest)
		return
	}

	subscriber := models.Subscriber{
		Email:    email,
		Timezone: offset,
	}
	err = m.DB.AddSubscriber(subscriber)
	if err != nil {
		if errors.Is(err, customerrors.DuplicatedKey) {
			http.Error(w, "Already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Successfully subscribed"})
}
