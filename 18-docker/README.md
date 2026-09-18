# Stage 18: Docker & Containers

Welcome to **Stage 18** of your Go mastery journey. One of Go's greatest strengths in modern cloud computing is its ability to compile into a single, self-contained, statically linked binary. This enables ultra-lightweight, secure container images measuring **under 15 MB** with zero OS baggage.

---

## 1. The Production Go Dockerfile Standard

A production-grade Go container must satisfy four non-negotiable rules:
1. **Multi-Stage Build**: Separate the heavy compilation environment (compilers, toolchains, git, caches) from the minimal execution environment.
2. **Static Compilation**: Set `CGO_ENABLED=0` so the compiled binary has no dynamic linking dependencies on host `libc` or `musl`.
3. **Strip Symbols**: Use `-ldflags="-s -w"` to strip DWARF debugging tables and symbol tables, reducing binary size by ~30-40%.
4. **Non-Root Execution**: Run under a dedicated unprivileged user (`nonroot:nonroot` or `UID 10001`) to protect host systems in case of container escapes.

```dockerfile
# Stage 1: Build Environment
FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata

# Cache dependency layer first
COPY go.mod go.sum ./
RUN go mod download

# Build statically linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/server .

# Stage 2: Minimal Runtime Environment
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/server /app/server

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/server"]
```

---

## 2. Directory Layout

```
18-docker/
├── README.md
├── examples/
│   ├── multistage_docker/         # Canonical multi-stage Dockerfile and Go server
│   └── compose_topology/          # Docker Compose with Go and PostgreSQL
├── exercises/
│   ├── 01-dockerfile-linter       # Go tool verifying container security rules
│   ├── 02-self-healthchecker      # Minimal binary self-health probe for distroless
│   └── 03-compose-config-validator# Docker Compose schema & security validator
├── solutions/
│   ├── 01-dockerfile-linter
│   ├── 02-self-healthchecker
│   └── 03-compose-config-validator
└── challenges/
    └── containerized-microservice # Complete microservice, Dockerfile, Compose, & tests
```
