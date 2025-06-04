package loadbalancer

import (
	"net"
	"net/http"
	"strconv"

	"github.com/PreethamVJ/LB/internal/config"
)

type Layer7 struct {
	*BaseLoadBalancer
	server *http.Server
}

func NewLayer7(config *config.Config) *Layer7 {
	return &Layer7{
		BaseLoadBalancer: NewBase(config),
	}
}

func (l *Layer7) Start() error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server, err := l.Algorithm.PickServer()
		if err != nil {
			http.Error(w, "No servers available", http.StatusServiceUnavailable)
			return
		}

		// Forward request (like Rust's Server::handle_request)
		if err := server.ForwardHTTP(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	})

	l.server = &http.Server{
		Addr:    net.JoinHostPort(l.Config.LoadBalancer.Address, strconv.Itoa(l.Config.LoadBalancer.Port)),
		Handler: handler,
	}

	return l.server.ListenAndServe()
}

func (l *Layer7) Stop() {
	_ = l.server.Close()
}
