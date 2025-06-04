package loadbalancer

import (
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/PreethamVJ/LB/internal/config"
)

type Layer4 struct {
	*BaseLoadBalancer
	listener net.Listener
	wg       sync.WaitGroup
	running  bool
}

func NewLayer4(config *config.Config) *Layer4 {
	return &Layer4{
		BaseLoadBalancer: NewBase(config),
	}
}

func (l *Layer4) Start() error {
	listener, err := net.Listen("tcp", net.JoinHostPort(
		l.Config.LoadBalancer.Address,
		strconv.Itoa(l.Config.LoadBalancer.Port),
	))
	if err != nil {
		return err
	}

	l.listener = listener
	l.running = true

	for l.running {
		conn, err := l.listener.Accept()
		if err != nil {
			if l.running {
				log.Printf("Accept error: %v", err)
			}
			continue
		}

		l.wg.Add(1)
		go l.handleConnection(conn)
	}

	return nil
}

func (l *Layer4) handleConnection(conn net.Conn) {
	defer l.wg.Done()
	defer conn.Close()

	server, err := l.Algorithm.PickServer()
	if err != nil {
		log.Printf("No available servers: %v", err)
		return
	}

	if err := server.TransferData(conn); err != nil {
		log.Printf("Transfer failed: %v", err)
	}
}

func (l *Layer4) Stop() {
	l.running = false
	if l.listener != nil {
		l.listener.Close()
	}
	l.wg.Wait()
}
