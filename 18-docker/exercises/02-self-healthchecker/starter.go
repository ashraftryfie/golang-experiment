package selfhealth

import (
	"errors"
	"time"
)

var (
	ErrNotImplemented  = errors.New("exercise not implemented yet")
	ErrUnhealthyStatus = errors.New("server returned unhealthy HTTP status")
)

const DefaultHealthURL = "http://127.0.0.1:8080/healthz"

func ParseCLIArgs(args []string) (isProbe bool, targetURL string) {
	return false, ""
}

func RunHealthProbe(targetURL string, timeout time.Duration) error {
	return ErrNotImplemented
}
