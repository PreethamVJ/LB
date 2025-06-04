package algorithm

import (
	"errors"
	"sync"

	"github.com/PreethamVJ/LB/server"
)

var ErrNoServers = errors.New("no servers available")

type RoundRobin struct {
	servers []*server.Server
	mu      sync.Mutex
	index   int
}

func NewRoundRobin(servers []*server.Server) *RoundRobin {
	return &RoundRobin{
		servers: servers,
		index:   0,
	}
}

func (rr *RoundRobin) PickServer() (*server.Server, error) {
	rr.mu.Lock()
	defer rr.mu.Unlock()

	if len(rr.servers) == 0 {
		return nil, ErrNoServers
	}

	server := rr.servers[rr.index]
	rr.index = (rr.index + 1) % len(rr.servers) // Cycle through servers

	return server, nil
}
