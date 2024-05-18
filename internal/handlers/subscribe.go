package handlers

import (
	"errors"
	"github.com/go-chi/render"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/forms"
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

	err := m.DB.AddSubscriber(email)
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
