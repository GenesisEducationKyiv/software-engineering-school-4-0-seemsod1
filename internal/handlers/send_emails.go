package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

type Notifier interface {
	SendRate(subscribers []string)
}

func (m *Repository) SendEmails(w http.ResponseWriter, _ *http.Request) {
	m.Logger.Info("Sending emails")

	subs, err := m.Subscriber.GetSubscribers()
	if err != nil {
		m.Logger.Error("Getting subscribers", zap.Error(err))
		http.Error(w, "Failed to get subscribers", http.StatusBadRequest)
		return
	}
	m.Notifier.SendRate(subs)

	w.WriteHeader(http.StatusOK)
}
