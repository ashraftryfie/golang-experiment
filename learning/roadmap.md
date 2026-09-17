# Comprehensive Go Learning Roadmap

This roadmap takes you from zero Go experience to designing, testing, and deploying a production microservice.

---

## The 4 Learning Epochs

```text
┌─────────────────────────────────────────────────────────┐
│ EPOCH 1: GO FUNDAMENTALS & TYPE MASTERY (Stages 00-06)  │
│ Toolchain, Types, Slices, Structs, Interfaces, Errors   │
└────────────────────────────┬────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────┐
│ EPOCH 2: STANDARD LIBRARY & IO PIPELINES (Stages 07-11) │
│ io.Reader/Writer, JSON, HTTP Routing, SQL & Persistence │
└────────────────────────────┬────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────┐
│ EPOCH 3: CONCURRENCY & SYSTEMS THINKING (Stages 12-15)  │
│ Table Tests, Goroutines, CSP Channels, Context Trees    │
└────────────────────────────┬────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────┐
│ EPOCH 4: PRODUCTION BACKENDS & DEVOPS (Stages 16-21)   │
│ Security, Middleware, SRE, Docker, Microservices, Cloud │
└─────────────────────────────────────────────────────────┘
```

---

## Detailed Stage Breakdown

### Epoch 1: Go Fundamentals & Type Mastery
- **Stage 00: Orientation & Environment**
  - Toolchain (`go run`, `go build`, `go vet`, `go fmt`), workspace layout, module initialization, compilation model.
- **Stage 01: Go Fundamentals**
  - Static typing, basic types, zero values, typed vs untyped constants, control flow (if with init, switch, for loop variations), functions and return values.
- **Stage 02: Data Structures**
  - Arrays vs slices, slice header mechanics (`ptr`, `len`, `cap`), `append()` and reallocation, sub-slicing gotchas, map internals and nil map traps.
- **Stage 03: Functions & Methods**
  - First-class functions, closures, variadic functions, `defer` lifecycle and evaluation timing, value vs pointer receivers.
- **Stage 04: Structs & Interfaces**
  - Struct definition, memory layout, embedding vs inheritance, implicit interface satisfaction, interface composition, empty interface / `any`, type assertions and switches.
- **Stage 05: Error Handling**
  - Errors as values, `error` interface, sentinel errors vs error types, wrapping with `fmt.Errorf("%w")`, `errors.Is` and `errors.As`, panic/recover philosophy.
- **Stage 06: Packages & Modules**
  - Package visibility (exported vs unexported), `internal/` directory enforcement, package naming conventions, cyclic dependency prevention, `go.mod` and `go.sum`.

### Epoch 2: Standard Library & I/O Pipelines
- **Stage 07: Standard Library Deep Dive**
  - `io.Reader` and `io.Writer` interfaces, stream processing, `bufio.Scanner` and `bufio.Writer`, string manipulation with `strings` and `bytes`, `time` duration math and tickers.
- **Stage 08: Files, JSON & Serialization**
  - `os` file handles and permissions, streaming `json.Decoder` vs `json.Unmarshal`, struct tags (`json:"..."`), implementing `json.Marshaler` and `json.Unmarshaler`.
- **Stage 09: HTTP Fundamentals**
  - Client-server architecture, `http.Handler` and `http.HandlerFunc`, Go 1.22+ method/path pattern matching in `http.ServeMux`, writing custom middlewares, query and body parsing.
- **Stage 10: REST API Engineering**
  - RESTful resource design, clean status codes, JSON error envelopes, request payload validation, route parameter extraction.
- **Stage 11: PostgreSQL & Persistence**
  - Relational database schema design, `database/sql` driver lifecycle, connection pool tuning (`SetMaxOpenConns`, `SetMaxIdleConns`), SQL migrations, atomic transactions (`BeginTx`).

### Epoch 3: Concurrency & Systems Thinking
- **Stage 12: Testing, Fuzzing & Benchmarking**
  - Idiomatic table-driven tests, subtests with `t.Run()`, test fixtures and golden files, `httptest` servers and recorders, race condition detection with `-race`, micro-benchmarks (`testing.B`), fuzz testing (`testing.F`).
- **Stage 13: Concurrency & Channels**
  - Go runtime scheduler (M:N scheduler, Goroutines vs OS threads), unbuffered vs buffered channels, channel ownership and closing rules, `select` statement with timeouts and default cases, `sync.WaitGroup`, `sync.Mutex` and `sync.RWMutex`, atomic operations.
- **Stage 14: Context & Cancellation**
  - `context.Context` hierarchy, cancellation propagation (`context.WithCancel`), deadlines and timeouts (`context.WithTimeout`), request-scoped values (`context.WithValue`) best practices and antipatterns.
- **Stage 15: Networking & Sockets**
  - Low-level TCP connections with `net.Dial` and `net.Listen`, line-delimited socket protocols, graceful connection draining, UDP basics.

### Epoch 4: Production Backends & Cloud Engineering
- **Stage 16: Authentication & Security**
  - Cryptographic hashing with `golang.org/x/crypto/bcrypt`, stateless JWT generation and validation, role-based access control (RBAC) middleware, timing attack mitigation, secure headers.
- **Stage 17: Production Backend Patterns**
  - Graceful shutdown patterns capturing `os.Interrupt` and `syscall.SIGTERM`, retry with exponential backoff and jitter, circuit breaker pattern, token-bucket rate limiting.
- **Stage 18: Docker & Containerization**
  - Multi-stage Docker builds, static CGO-free compilation, minimal non-root distroless runtime containers, Docker Compose for multi-container development (App + Postgres + Redis).
- **Stage 19: Observability & Telemetry**
  - Structured logging with `log/slog` (text vs JSON handlers), correlation IDs / trace IDs across middlewares, Prometheus application metrics (counters, histograms), OpenTelemetry tracing foundations.
- **Stage 20: Microservice Architecture**
  - Hexagonal / Ports & Adapters architecture, separation of domain models from transport and persistence layers, interface boundaries, asynchronous background task workers.
- **Stage 21: Capstone Integration Service**
  - Comprehensive production microservice connecting HTTP REST, PostgreSQL storage, Redis caching, worker queue processing, simulated external API integration, resilience, and full automated test suite.
