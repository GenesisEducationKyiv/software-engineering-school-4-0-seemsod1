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

type Node struct {
	rateService RateService
	next        Chain
}

func NewNode(fetcher RateService) *Node {
	return &Node{
		rateService: fetcher,
	}
}

func (n *Node) SetNext(chainInterface Chain) {
	n.next = chainInterface
}

func (n *Node) GetRate(ctx context.Context, base, target string) (float64, error) {
	rate, err := n.rateService.GetRate(ctx, base, target)
	if err != nil {
		next := n.next
		if next == nil {
			return -1, ErrNoRateProviders
		}

		return next.GetRate(ctx, base, target)
	}
	return rate, nil
}
