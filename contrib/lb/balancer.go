package lb

import (
	"errors"
	"github.com/zander-84/seagull/contract"
)

// ErrNoNode is returned when no qualifying node are available.
var ErrNoNode = errors.New("no node available")
var ErrNotImplemented = errors.New("not implemented")

type Policy int

const (
	RoundRobin Policy = iota
	WeightRoundRobin
	ConsistentHash
	Random
)

func NewBalancer(listener Listener, p Policy, record bool) contract.Balancer {
	switch p {
	case RoundRobin:
		return NewRoundRobin(listener, record)
	case WeightRoundRobin:
		return NewWeightRoundRobin(listener, record)
	case ConsistentHash:
		return NewConsistentHash(listener, record)
	case Random:
		return NewRandom(listener, record)
	default:
		return NewRoundRobin(listener, record)
	}
}
