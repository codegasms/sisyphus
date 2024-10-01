package main

import (
	"errors"
	"math/rand/v2"
)

type StrategyKind string

const (
	StrategyRandom             StrategyKind = "random"
	StrategyRoundRobin                      = "round-robin"
	StrategyWeightedRoundRobin              = "weighted-round-robin"
	StrategyLeastConnections                = "least-connections"
)

func StrategyFromConfig(config *Config) (Strategy, error) {
	switch config.Strategy {
	case StrategyRandom:
		return NewRandomStrategy(config.Servers), nil
	case StrategyRoundRobin:
		return NewRoundRobinStrategy(config.Servers), nil
	case StrategyWeightedRoundRobin:
		if len(config.Servers) != len(config.Weights) {
			return nil, errors.New("weights empty")
		}
		return NewWeightedRoundRobinStrategy(config.Servers, config.Weights), nil
	case StrategyLeastConnections:
		return NewLeastConnectionsStrategy(config.Servers), nil
	}

	return nil, errors.New("invalid strategy in config")
}

type ServerAddr string

type Strategy interface {
	ServerAddr() (ServerAddr, error)
}

type RandomStrategy struct {
	servers []ServerAddr
}

func NewRandomStrategy(servers []ServerAddr) *RandomStrategy {
	return &RandomStrategy{servers}
}

func (strategy *RandomStrategy) ServerAddr() (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	server := strategy.servers[rand.IntN(len(strategy.servers))]

	return server, nil
}

type RoundRobinStrategy struct {
	servers   []ServerAddr
	nextIndex int
}

func NewRoundRobinStrategy(servers []ServerAddr) *RoundRobinStrategy {
	return &RoundRobinStrategy{servers, 0}
}

func (strategy *RoundRobinStrategy) ServerAddr() (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	server := strategy.servers[strategy.nextIndex]
	strategy.nextIndex = (strategy.nextIndex + 1) % len(strategy.servers)

	return server, nil
}

type WeightedRoundRobinStrategy struct {
	servers       []ServerAddr
	weights       []float32
	nextIndex     int
	selectedTimes int
}

func NewWeightedRoundRobinStrategy(servers []ServerAddr, weights []float32) *WeightedRoundRobinStrategy {
	return &WeightedRoundRobinStrategy{servers, weights, 0, 0}
}

func (strategy *WeightedRoundRobinStrategy) ServerAddr() (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	server := strategy.servers[strategy.nextIndex]
	strategy.selectedTimes++

	// TODO: Implement a better weight conversion strategy
	if strategy.selectedTimes >= int(strategy.weights[strategy.nextIndex]*100) {
		strategy.nextIndex = (strategy.nextIndex + 1) % len(strategy.servers)
		strategy.selectedTimes = 0
	}

	return server, nil
}

type LeastConnectionsStrategy struct {
	servers     []ServerAddr
	connections []int
}

func NewLeastConnectionsStrategy(servers []ServerAddr) *LeastConnectionsStrategy {
	connections := make([]int, len(servers))
	for i := range connections {
		connections[i] = 0
	}

	return &LeastConnectionsStrategy{servers, connections}
}

func (strategy *LeastConnectionsStrategy) ServerAddr() (ServerAddr, error) {
	return "", nil
}
