package rateapi

import (
	"log"
)

type RateService interface {
	GetRate(base, target string) (float64, error)
}

// LoggingClient is a client that logs the responses
type LoggingClient struct {
	name        string
	rateService RateService
}

func NewLoggingClient(name string, rateService RateService) *LoggingClient {
	return &LoggingClient{
		name:        name,
		rateService: rateService,
	}
}

func (l *LoggingClient) GetRate(base, target string) (float64, error) {
	rate, err := l.rateService.GetRate(base, target)
	if err != nil {
		log.Printf("%s: Response: {error: %v}", l.name, err)
		return -1, err
	}

	log.Printf("%s: Response: {price: %f}", l.name, rate)
	return rate, nil
}
