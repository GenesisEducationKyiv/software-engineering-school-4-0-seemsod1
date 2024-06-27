package chain

import (
	"context"
	"errors"
)

type RateService interface {
	GetRate(ctx context.Context, base, target string) (float64, error)
}

type Chain interface {
	RateService
	SetNext(chainInterface Chain)
}

var ErrNoRateProviders = errors.New("no available providers to get rate")

type BaseChain struct {
	rateService RateService
	next        Chain
}

func NewBaseChain(fetcher RateService) *BaseChain {
	return &BaseChain{rateService: fetcher}
}

func (b *BaseChain) SetNext(chainInterface Chain) {
	b.next = chainInterface
}

func (b *BaseChain) GetRate(ctx context.Context, base, target string) (float64, error) {
	rate, err := b.rateService.GetRate(ctx, base, target)
	if err != nil {
		next := b.next
		if next == nil {
			return -1, ErrNoRateProviders
		}

		return next.GetRate(ctx, base, target)
	}
	return rate, nil
}
