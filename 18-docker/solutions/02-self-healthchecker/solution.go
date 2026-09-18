package selfhealth

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrUnhealthyStatus = errors.New("server returned unhealthy HTTP status")
)

const DefaultHealthURL = "http://127.0.0.1:8080/healthz"

func ParseCLIArgs(args []string) (isProbe bool, targetURL string) {
	for _, arg := range args {
		if arg == "-healthcheck" || arg == "--healthcheck" {
			return true, DefaultHealthURL
		}
		if strings.HasPrefix(arg, "-healthcheck=") {
			return true, strings.TrimPrefix(arg, "-healthcheck=")
		}
		if strings.HasPrefix(arg, "--healthcheck=") {
			return true, strings.TrimPrefix(arg, "--healthcheck=")
		}
	}
	return false, ""
}

func RunHealthProbe(targetURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return fmt.Errorf("failed creating health request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("health probe network error: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("%w: %d", ErrUnhealthyStatus, resp.StatusCode)
	}

	return nil
}
