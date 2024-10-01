package main

type StrategyKind int

const (
	StrategyRandom StrategyKind = iota
	StrategyRoundRobin
	StrategyWeightedRoundRobin
	StrategyLeastConnections
)

var strategyName = map[StrategyKind]string{
	StrategyRandom:             "random",
	StrategyRoundRobin:         "round-robin",
	StrategyWeightedRoundRobin: "weighted-round-robin",
	StrategyLeastConnections:   "least-connections",
}

func (s StrategyKind) String() string {
	return strategyName[s]
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

type RoundRobinStrategy struct {
	servers   []ServerAddr
	nextIndex int
}

func NewRoundRobinStrategy(servers []ServerAddr) *RoundRobinStrategy {
	return &RoundRobinStrategy{servers, 0}
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
