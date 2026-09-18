package gracefulserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	ErrServerNotRunning = errors.New("server is not running")
)

type GracefulServer struct {
	addr     string
	handler  http.Handler
	listener net.Listener
	server   *http.Server

	ctx        context.Context
	cancelFunc context.CancelFunc
	taskWG     sync.WaitGroup

	mu      sync.Mutex
	running bool
}

func NewGracefulServer(addr string, handler http.Handler) *GracefulServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &GracefulServer{
		addr:       addr,
		handler:    handler,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

func (s *GracefulServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed listening on %s: %w", s.addr, err)
	}
	s.listener = l
	s.server = &http.Server{
		Handler: s.handler,
	}
	s.running = true

	go func() {
		_ = s.server.Serve(s.listener)
	}()

	return nil
}

func (s *GracefulServer) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

// TrackTask executes a background task with the server's lifecycle context, tracked until completion.
func (s *GracefulServer) TrackTask(fn func(ctx context.Context)) {
	s.taskWG.Add(1)
	go func() {
		defer s.taskWG.Done()
		fn(s.ctx)
	}()
}

// Shutdown drains in-flight requests and background workers within the specified timeout.
func (s *GracefulServer) Shutdown(timeout time.Duration) error {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return ErrServerNotRunning
	}
	s.running = false
	s.mu.Unlock()

	// 1. Signal background tasks to stop
	s.cancelFunc()

	// 2. Bound shutdown duration
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Channel to signal completion of both server shutdown and task drain
	doneCh := make(chan error, 1)

	go func() {
		var err error
		if s.server != nil {
			err = s.server.Shutdown(shutdownCtx)
		}
		// Wait for background tasks
		s.taskWG.Wait()
		doneCh <- err
	}()

	select {
	case err := <-doneCh:
		return err
	case <-shutdownCtx.Done():
		return shutdownCtx.Err()
	}
}
