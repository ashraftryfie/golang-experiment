# Exercise 01: Go Dockerfile Linter

## Objective
Implement a Go static analysis linter for Dockerfiles that enforces cloud-native best practices, layer caching, and container security.

## Requirements
1. **Multi-Stage Build**: Ensure at least two `FROM` instructions exist (`MULTI_STAGE_REQUIRED`).
2. **CGO Disabled**: Ensure `CGO_ENABLED=0` is present during compilation to prevent dynamic libc dependencies (`CGO_DISABLED_REQUIRED`).
3. **Dependency Layer Caching**: Ensure `COPY go.mod` or `COPY go.sum` occurs before `COPY . .` to maximize Docker build caching (`DEPENDENCY_CACHE_ORDER`).
4. **Non-Root Execution**: Ensure the final stage defines an explicit non-root user (`USER nonroot` or `USER 10001` or `USER nonroot:nonroot`) before the entrypoint (`NON_ROOT_USER_REQUIRED`).

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./18-docker/exercises/01-dockerfile-linter`
