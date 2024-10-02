package main

import (
	"errors"
	"log"
	"math/rand/v2"
	"strings"
)

type StrategyKind string

const (
	StrategyRandom                   StrategyKind = "random"
	StrategyRoundRobin                            = "round-robin"
	StrategyWeightedRoundRobin                    = "weighted-round-robin"
	StrategyIPHash                                = "ip-hash"
	StrategyLeastConnections                      = "least-connections"
	StrategyweightedLeastConnections              = "weighted-least-connections"
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
	case StrategyIPHash:
		return NewIPHashStrategy(config.Servers), nil
	case StrategyLeastConnections:
		return NewLeastConnectionsStrategy(config.Servers), nil
	case StrategyweightedLeastConnections:
		if len(config.Servers) != len(config.Weights) {
			return nil, errors.New("weights empty")
		}
		return WeightedNewLeastConnectionsStrategy(config.Servers, config.Weights), nil
	}

	return nil, errors.New("invalid strategy in config")
}

type ServerAddr string

type Strategy interface {
	ServerAddr(ClientAddr string) (ServerAddr, error)

	// TODO: Methods to notify the strategy about connection opening and closing.
	Connected(addr ServerAddr)
	Disconnected(addr ServerAddr)
}

const WeightFactor float32 = 100 // Weight Factor for WeightedRoundRobinStrategy

type RandomStrategy struct {
	servers []ServerAddr
}

func NewRandomStrategy(servers []ServerAddr) *RandomStrategy {
	return &RandomStrategy{servers}
}

func (strategy *RandomStrategy) ServerAddr(ClientAddr string) (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	server := strategy.servers[rand.IntN(len(strategy.servers))]

	return server, nil
}

func (strategy *RandomStrategy) Connected(addr ServerAddr)    {}
func (strategy *RandomStrategy) Disconnected(addr ServerAddr) {}

type RoundRobinStrategy struct {
	servers   []ServerAddr
	nextIndex int
}

func NewRoundRobinStrategy(servers []ServerAddr) *RoundRobinStrategy {
	return &RoundRobinStrategy{servers, 0}
}

func (strategy *RoundRobinStrategy) ServerAddr(ClientAddr string) (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	server := strategy.servers[strategy.nextIndex]
	strategy.nextIndex = (strategy.nextIndex + 1) % len(strategy.servers)

	return server, nil
}

func (strategy *RoundRobinStrategy) Connected(addr ServerAddr)    {}
func (strategy *RoundRobinStrategy) Disconnected(addr ServerAddr) {}

type WeightedRoundRobinStrategy struct {
	servers       []ServerAddr
	weights       []float32
	nextIndex     int
	selectedTimes int
}

func NewWeightedRoundRobinStrategy(servers []ServerAddr, weights []float32) *WeightedRoundRobinStrategy {
	// Normalize the weights to 1
	sum := float32(0)
	for _, weight := range weights {
		sum += weight
	}

	for i, weight := range weights {
		weights[i] = weight / sum
	}

	return &WeightedRoundRobinStrategy{servers, weights, 0, 0}
}

func (strategy *WeightedRoundRobinStrategy) ServerAddr(ClientAddr string) (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	server := strategy.servers[strategy.nextIndex]
	strategy.selectedTimes++

	if strategy.selectedTimes >= int(strategy.weights[strategy.nextIndex]*WeightFactor) {
		strategy.nextIndex = (strategy.nextIndex + 1) % len(strategy.servers)
		strategy.selectedTimes = 0
	}

	return server, nil
}

func (strategy *WeightedRoundRobinStrategy) Connected(addr ServerAddr)    {}
func (strategy *WeightedRoundRobinStrategy) Disconnected(addr ServerAddr) {}

type IPHashStrategy struct {
	servers []ServerAddr
}

func NewIPHashStrategy(servers []ServerAddr) *IPHashStrategy {
	return &IPHashStrategy{servers}
}

func (strategy *IPHashStrategy) ServerAddr(ClientAddr string) (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	fragments := strings.Split(ClientAddr, ":")
	addr := strings.Join(fragments[:len(fragments)-1], ":")

	// Hash the IP address (Based on the Sum of Bytes of the Client Address) to select a server.
	// TODO: Implement a better hashing algorithm.
	hash := 0
	for i := 0; i < len(addr); i++ {
		hash += int(addr[i])
	}

	server := strategy.servers[hash%len(strategy.servers)]

	return server, nil
}

func (strategy *IPHashStrategy) Connected(addr ServerAddr)    {}
func (strategy *IPHashStrategy) Disconnected(addr ServerAddr) {}

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

func (strategy *LeastConnectionsStrategy) ServerAddr(ClientAddr string) (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	minConnections := strategy.connections[0]
	minIndex := 0

	for i, connections := range strategy.connections {
		if connections < minConnections {
			minConnections = connections
			minIndex = i
		}
	}

	server := strategy.servers[minIndex]
	strategy.connections[minIndex]++

	return server, nil
}

func (strategy *LeastConnectionsStrategy) Connected(addr ServerAddr) {}

func (strategy *LeastConnectionsStrategy) Disconnected(addr ServerAddr) {
	log.Println(strategy.connections)
	for i, serverAddr := range strategy.servers {
		if serverAddr == addr && strategy.connections[i] > 0 {
			strategy.connections[i]--
			break
		}
	}
	log.Println(strategy.connections)
}

type WeightedLeastConnectionsStrategy struct {
	servers     []ServerAddr
	connections []int
	weights     []float32
}

func WeightedNewLeastConnectionsStrategy(servers []ServerAddr, weights []float32) *WeightedLeastConnectionsStrategy {
	connections := make([]int, len(servers))
	for i := range connections {
		connections[i] = 0
	}

	// Normalize the weights to 1
	sum := float32(0)
	for _, weight := range weights {
		sum += weight
	}

	for i, weight := range weights {
		weights[i] = weight / sum
	}

	return &WeightedLeastConnectionsStrategy{servers, connections, weights}
}

func (strategy *WeightedLeastConnectionsStrategy) ServerAddr(ClientAddr string) (ServerAddr, error) {
	if len(strategy.servers) == 0 {
		return "", errors.New("no servers available")
	}

	maxRatio := float32(0.0)
	maxIndex := 0

	totalConnections := 0
	for _, connections := range strategy.connections {
		totalConnections += connections
	}

	// Below code works for Least Connection Strategy
	for i, connections := range strategy.connections {
		normalizedConnections := float32(connections) / (float32(totalConnections) + 1e-6)

		ratio := float32(strategy.weights[i]) / (float32(normalizedConnections) + 1e-6)

		if ratio > maxRatio {
			maxRatio = ratio
			maxIndex = i
		}
	}

	server := strategy.servers[maxIndex]
	strategy.connections[maxIndex]++

	return server, nil
}

func (strategy *WeightedLeastConnectionsStrategy) Connected(addr ServerAddr) {}

func (strategy *WeightedLeastConnectionsStrategy) Disconnected(addr ServerAddr) {
	log.Println(strategy.connections)
	for i, serverAddr := range strategy.servers {
		if serverAddr == addr && strategy.connections[i] > 0 {
			strategy.connections[i]--
			break
		}
	}
	log.Println(strategy.connections)
}
