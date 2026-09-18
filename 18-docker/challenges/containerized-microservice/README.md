# Engineering Challenge: Containerized Production Microservice

## Challenge Overview
Deploying Go services in modern cloud architectures (Kubernetes, AWS ECS, Google Cloud Run) demands strict compliance with container standards: tiny immutable images, non-root execution, dual liveness/readiness probes, and self-contained health checking.

You must build a complete **Containerized Production Microservice Workspace** comprising:
1. **Cloud-Native Go Server (`server.go`)**:
   - `/healthz`: Kubernetes/Docker liveness probe (returns `200 OK`).
   - `/readyz`: Readiness probe reflecting external dependency health (e.g. database connectivity). Returns `200 OK` when healthy and `503 Service Unavailable` when degraded.
   - Built-in self-healthcheck probe triggered via `-healthcheck` flag for distroless compatibility.
2. **Production Multi-Stage Dockerfile (`Dockerfile`)**:
   - Compiles statically with `CGO_ENABLED=0` and `-ldflags="-s -w"`.
   - Leverages `go.mod` layer caching.
   - Packages into `gcr.io/distroless/static-debian12:nonroot` running as an unprivileged user.
3. **Docker Compose Topology (`docker-compose.yml`)**:
   - Connects the application to a PostgreSQL database.
   - Enforces database healthchecks and depends-on conditions.

---

## Running Tests
Run tests with race detection:
```bash
go test -v -race ./18-docker/challenges/containerized-microservice/...
```
