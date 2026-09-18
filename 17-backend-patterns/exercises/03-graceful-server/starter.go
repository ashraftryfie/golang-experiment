package gracefulserver

import (
	"context"
	"errors"
	"net/http"
	"time"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

type GracefulServer struct{}

func NewGracefulServer(addr string, handler http.Handler) *GracefulServer {
	return nil
}

func (s *GracefulServer) Start() error {
	return ErrNotImplemented
}

func (s *GracefulServer) Addr() string {
	return ""
}

func (s *GracefulServer) TrackTask(fn func(ctx context.Context)) {
}

func (s *GracefulServer) Shutdown(timeout time.Duration) error {
	return ErrNotImplemented
}
