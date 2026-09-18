# Exercise 02: Binary Self-Health Probe for Distroless

## Objective
Implement an internal health probe within a Go binary so Docker and Kubernetes can execute healthchecks on minimal `scratch` or `distroless` containers without needing external utilities like `curl` or `wget`.

## Requirements
1. **`ParseCLIArgs(args []string) (isProbe bool, targetURL string)`**:
   - Detects flags like `-healthcheck`, `--healthcheck`, or `-healthcheck=http://...`.
   - Defaults to `http://127.0.0.1:8080/healthz` if no custom URL is specified.
2. **`RunHealthProbe(targetURL string, timeout time.Duration) error`**:
   - Performs an HTTP GET with context timeout.
   - Consumes and closes response body to prevent socket leaks.
   - Returns `nil` if status is in the 2xx range; returns `ErrUnhealthyStatus` or network errors otherwise.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./18-docker/exercises/02-self-healthchecker`
