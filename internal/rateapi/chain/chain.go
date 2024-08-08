package chain

import (
	"context"
	"errors"
)

type rateService interface {
	GetRate(ctx context.Context, base, target string) (float64, error)
}

type Chain interface {
	rateService
	SetNext(chainInterface Chain)
}

var ErrNoRateProviders = errors.New("no available providers to get rate")

type Node struct {
	rateService rateService
	next        Chain
}

func NewNode(fetcher rateService) *Node {
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
