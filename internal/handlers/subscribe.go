package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/render"
	customerrors "github.com/seemsod1/api-project/internal/errors"
	"github.com/seemsod1/api-project/internal/forms"
	"github.com/seemsod1/api-project/internal/models"
	"github.com/seemsod1/api-project/internal/timezone"
)

// Subscribe subscribes a user to the newsletter
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(fmt.Errorf("parsing form: %w", err))
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("email")
	form.IsEmail("email")

	if !form.Valid() {
		log.Println("Invalid email")
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	offset, err := timezone.ProcessTimezoneHeader(r)
	if err != nil {
		http.Error(w, "Invalid timezone", http.StatusBadRequest)
		return
	}

	if err = m.DB.AddSubscriber(models.Subscriber{
		Email:    email,
		Timezone: offset,
	}); err != nil {
		if errors.Is(err, customerrors.ErrDuplicatedKey) {
			log.Printf("%s: %s", email, "Already exists")
			http.Error(w, "Already exists", http.StatusConflict)
			return
		}
		log.Println(fmt.Errorf("adding subscriber: %w", err))
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "Successfully subscribed"})
}
