package handlers

import (
	"log"
	"net/http"
)

func (m *Repository) Rate(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Rate endpoint")); err != nil {
		log.Fatal(err)
	}
}
