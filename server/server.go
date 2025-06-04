package server

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Address        string
	Port           int
	MaxConnections int
	Weight         int

	mu                  sync.RWMutex
	currentConnections  int
	totalConnections    uint64
	successfulTransfers uint64
	failedTransfers     uint64
	lastRequestTime     time.Time
	responseTime        time.Duration
	isAlive             bool
}

// TransferData handles bidirectional raw TCP proxying between clientConn and backend server.
func (s *Server) TransferData(clientConn net.Conn) error {
	// Track connection count
	s.mu.Lock()
	s.currentConnections++
	s.totalConnections++
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.currentConnections--
		s.mu.Unlock()
	}()

	// Connect to backend server
	serverConn, err := net.Dial("tcp", net.JoinHostPort(s.Address, strconv.Itoa(s.Port)))
	if err != nil {
		s.mu.Lock()
		s.failedTransfers++
		s.mu.Unlock()
		return err
	}
	defer serverConn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	var copyErr error
	var copyErrMutex sync.Mutex

	// Helper to record first error
	recordError := func(err error) {
		if err != nil {
			copyErrMutex.Lock()
			if copyErr == nil {
				copyErr = err
			}
			copyErrMutex.Unlock()
		}
	}

	go func() {
		defer wg.Done()
		_, err := io.Copy(serverConn, clientConn)
		recordError(err)
		serverConn.Close() // close to unblock other copy
	}()

	go func() {
		defer wg.Done()
		_, err := io.Copy(clientConn, serverConn)
		recordError(err)
		clientConn.Close() // close to unblock other copy
	}()

	wg.Wait()

	s.mu.Lock()
	if copyErr == nil {
		s.successfulTransfers++
	} else {
		s.failedTransfers++
	}
	s.mu.Unlock()

	return copyErr
}

// ForwardHTTP proxies HTTP requests from the client to the backend server.
func (s *Server) ForwardHTTP(w http.ResponseWriter, r *http.Request) error {
	// Update X-Forwarded-For header to append client IP
	clientIP := r.RemoteAddr
	if ipPort := strings.Split(clientIP, ":"); len(ipPort) > 0 {
		clientIP = ipPort[0]
	}

	// Append to existing X-Forwarded-For or set if empty
	if prior, ok := r.Header["X-Forwarded-For"]; ok {
		r.Header.Set("X-Forwarded-For", strings.Join(prior, ", ")+", "+clientIP)
	} else {
		r.Header.Set("X-Forwarded-For", clientIP)
	}

	// Set X-Forwarded-Host to original host
	r.Header.Set("X-Forwarded-Host", r.Host)

	// Adjust URL to point to backend server
	r.URL.Scheme = "http" // Adjust if backend is https
	r.URL.Host = net.JoinHostPort(s.Address, strconv.Itoa(s.Port))

	// Send request to backend
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		s.mu.Lock()
		s.failedTransfers++
		s.mu.Unlock()
		return err
	}
	defer resp.Body.Close()

	// Copy headers to client response
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		s.mu.Lock()
		s.failedTransfers++
		s.mu.Unlock()
		return err
	}

	s.mu.Lock()
	s.successfulTransfers++
	s.mu.Unlock()
	return nil
}
