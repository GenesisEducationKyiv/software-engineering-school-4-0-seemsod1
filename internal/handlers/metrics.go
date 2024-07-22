package handlers

import (
	"github.com/VictoriaMetrics/metrics"
	"net/http"
)

func Metrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	metrics.WritePrometheus(w, false)
}
