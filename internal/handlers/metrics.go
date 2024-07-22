package handlers

import (
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

func Metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	metrics.WritePrometheus(w, false)
}
