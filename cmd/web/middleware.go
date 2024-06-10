package main

import "net/http"

// EnableCORS enables CORS
func EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {

			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,Set-Cookie,Date,Accept, Content-Type, X-CSRF-Token")
			return
		}
		h.ServeHTTP(w, r)
	})
}
