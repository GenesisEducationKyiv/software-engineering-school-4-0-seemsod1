package middlewarepkg

import (
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

// EnableCORS enables CORS
func EnableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Set-Cookie, Date, Accept, Content-Type, X-CSRF-Token")
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func Handle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := fmt.Sprintf(`requests_total{path=%q}`, r.URL.Path)
		metrics.GetOrCreateCounter(s).Inc()
		h.ServeHTTP(w, r)
	})
}
