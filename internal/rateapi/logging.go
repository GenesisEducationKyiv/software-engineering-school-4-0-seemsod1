package rateapi

import (
	"context"

	"github.com/seemsod1/api-project/pkg/logger"
)

type RateService interface {
	GetRate(ctx context.Context, base, target string) (float64, error)
}

// LoggingClient is a client that logs the responses
type LoggingClient struct {
	name        string
	rateService RateService
	logger      *logger.Logger
}

func NewLoggingClient(name string, rateService RateService, log *logger.Logger) *LoggingClient {
	return &LoggingClient{
		name:        name,
		rateService: rateService,
		logger:      log,
	}
}

func (l *LoggingClient) GetRate(ctx context.Context, base, target string) (float64, error) {
	rate, err := l.rateService.GetRate(ctx, base, target)
	if err != nil {
		l.logger.Warnf("%s: Response: {error: %v}", l.name, err)
		return -1, err
	}

	l.logger.Infof("%s: Response: {price: %f}", l.name, rate)
	return rate, nil
}
