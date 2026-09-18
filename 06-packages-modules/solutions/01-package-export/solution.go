package config

import (
	"errors"
	"time"
)

var (
	ErrEmptyHost      = errors.New("host cannot be empty")
	ErrInvalidPort    = errors.New("port must be between 1 and 65535")
	ErrInvalidTimeout = errors.New("timeout must be greater than zero")
)

type ServerConfig struct {
	host       string
	port       int
	timeoutSec int
}

func NewServerConfig(host string, port, timeoutSec int) (*ServerConfig, error) {
	if host == "" {
		return nil, ErrEmptyHost
	}
	if port < 1 || port > 65535 {
		return nil, ErrInvalidPort
	}
	if timeoutSec <= 0 {
		return nil, ErrInvalidTimeout
	}
	return &ServerConfig{
		host:       host,
		port:       port,
		timeoutSec: timeoutSec,
	}, nil
}

func (c *ServerConfig) Host() string {
	return c.host
}

func (c *ServerConfig) Port() int {
	return c.port
}

func (c *ServerConfig) Timeout() time.Duration {
	return time.Duration(c.timeoutSec) * time.Second
}
