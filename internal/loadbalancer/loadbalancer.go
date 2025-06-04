package loadbalancer

import (
	"github.com/PreethamVJ/LB/internal/config"
	"github.com/PreethamVJ/LB/internal/loadbalancer/algorithm"
	"github.com/PreethamVJ/LB/server"
)

type Algorithm interface {
	PickServer() (*server.Server, error)
}

type LoadBalancer interface {
	Start() error
	Stop()
}

type BaseLoadBalancer struct {
	Config    *config.Config
	Servers   []*server.Server
	Algorithm Algorithm
}

func NewBase(config *config.Config) *BaseLoadBalancer {
	servers := make([]*server.Server, 0, len(config.LoadBalancer.Servers))
	for _, s := range config.LoadBalancer.Servers {
		servers = append(servers, &server.Server{
			Address:        s.Address,
			Port:           s.Port,
			MaxConnections: s.MaxConnections,
			Weight:         s.Weight,
			// isAlive:        true, // Cannot set unexported field
		})
		// If server.Server has a method to set isAlive, call it here, e.g.:
		// servers[len(servers)-1].SetAlive(true)
	}
	var algo Algorithm
	switch config.LoadBalancer.Algorithm {
	case "round_robin":
		algo = algorithm.NewRoundRobin(servers)
	// Add other algorithms (least_conn, etc.)
	default:
		algo = algorithm.NewRoundRobin(servers) // Default
	}

	return &BaseLoadBalancer{
		Config:    config,
		Servers:   servers,
		Algorithm: algo,
	}
}
