package config

import (
	"errors"
	"time"
)

var (
	ErrEmptyHost      = errors.New("host cannot be empty")
	ErrInvalidPort    = errors.New("port must be between 1 and 65535")
	ErrInvalidTimeout = errors.New("timeout must be greater than zero")
	ErrNotImplemented = errors.New("TODO: implement")
)

type ServerConfig struct {
	// TODO: Define unexported fields for host, port, timeoutSec
}

func NewServerConfig(host string, port, timeoutSec int) (*ServerConfig, error) {
	// TODO: Validate inputs and return initialized *ServerConfig
	return nil, ErrNotImplemented
}

func (c *ServerConfig) Host() string {
	return ""
}

func (c *ServerConfig) Port() int {
	return 0
}

func (c *ServerConfig) Timeout() time.Duration {
	return 0
}
