# GoLang-Experiment

> An authentic, test-driven engineering journey from absolute Go fundamentals to a production-grade microservice backend.

[![CI](https://github.com/ashraftryfie/golang-experiment/actions/workflows/ci.yml/badge.svg)](https://github.com/ashraftryfie/golang-experiment/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ashraftryfie/golang-experiment)](https://goreportcard.com/report/github.com/ashraftryfie/golang-experiment)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)

---

## 🧭 Repository Mission

This repository is built for three simultaneous purposes:
1. **Mastery Learning**: Developing a deep, first-principles understanding of Go's type system, runtime, memory layout, concurrency primitives, and idiomatic backend patterns.
2. **Long-Term Knowledge Base**: A memory-first revision system allowing quick recall months later via structured cheat sheets, mental models, and "Remember it like this" anchors.
3. **Engineering Portfolio**: Demonstrating authentic, test-validated software engineering with zero fluff, genuine Git commits, clean architecture, automated CI, and production-ready microservice code.

---

## 📊 Learning Progress Dashboard

```text
Progress: [██░░░░░░░░░░░░░░░░░░] 5%
Verified Milestones: 1 / 22 Stages
Active Stage: 01 - Go Fundamentals
```

*Authoritative progress is tracked with test verification evidence in [`learning/progress.yaml`](./learning/progress.yaml).*

---

## 🗺️ Curriculum Roadmap

| Stage | Topic | Status | Exercises | Challenge | Core Competency |
| :---: | :--- | :---: | :---: | :---: | :--- |
| **00** | [Orientation & Environment](./00-orientation/README.md) | `MASTERED` | 1/1 | Passed | Toolchain, compilation, GC basics, formatting |
| **01** | [Go Fundamentals](./01-fundamentals/README.md) | `IN_PROGRESS` | 0/3 | Pending | Types, zero values, control flow, functions, testing |
| **02** | [Data Structures](./02-data-structures/README.md) | `IN_PROGRESS` | 0/3 | Pending | Arrays, slices, slice headers, maps, capacity |
| **03** | [Functions & Methods](./03-functions-methods/README.md) | `NOT_STARTED` | 0/3 | Pending | First-class funcs, closures, defer, value vs pointer receivers |
| **04** | [Structs & Interfaces](./04-structs-interfaces/README.md) | `NOT_STARTED` | 0/3 | Pending | Composition, embedding, implicit interfaces, polymorphism |
| **05** | [Error Handling](./05-error-handling/README.md) | `NOT_STARTED` | 0/3 | Pending | Errors as values, wrapping `%w`, `errors.Is`, `errors.As` |
| **06** | [Packages & Modules](./06-packages-modules/README.md) | `NOT_STARTED` | 0/3 | Pending | Module cache, `internal/`, package boundaries, visibility |
| **07** | [Standard Library Deep Dive](./07-standard-library/README.md) | `NOT_STARTED` | 0/3 | Pending | `io.Reader/Writer`, `bufio`, `time`, `bytes`, `strings` |
| **08** | [Files, JSON & Serialization](./08-files-json/README.md) | `NOT_STARTED` | 0/3 | Pending | Struct tags, streaming encoders, custom Marshaler, file I/O |
| **09** | [HTTP Fundamentals](./09-http/README.md) | `NOT_STARTED` | 0/3 | Pending | `http.Handler`, `ServeMux` (Go 1.22+ routing), middlewares |
| **10** | [REST API Engineering](./10-rest-api/README.md) | `NOT_STARTED` | 0/3 | Pending | REST principles, request validation, error responses |
| **11** | [PostgreSQL & Storage](./11-postgresql/README.md) | `NOT_STARTED` | 0/3 | Pending | `database/sql`, connection pooling, migrations, transactions |
| **12** | [Testing & Benchmarking](./12-testing/README.md) | `NOT_STARTED` | 0/3 | Pending | Table tests, httptest, race detector, benchmarks, fuzzing |
| **13** | [Concurrency & Channels](./13-concurrency/README.md) | `NOT_STARTED` | 0/3 | Pending | Goroutines, CSP channels, select, worker pools, sync.Mutex |
| **14** | [Context & Cancellation](./14-context/README.md) | `NOT_STARTED` | 0/3 | Pending | Cancellation trees, timeouts, dead-lines, request tracing |
| **15** | [Networking & Sockets](./15-networking/README.md) | `NOT_STARTED` | 0/3 | Pending | TCP listeners, UDP, graceful socket shutdown, keep-alives |
| **16** | [Authentication & Security](./16-auth/README.md) | `NOT_STARTED` | 0/3 | Pending | Password hashing (bcrypt), JWT, RBAC middleware, secrets |
| **17** | [Production Backend Patterns](./17-backend-patterns/README.md) | `NOT_STARTED` | 0/3 | Pending | Graceful shutdown, circuit breakers, rate limiting, retries |
| **18** | [Docker & Containers](./18-docker/README.md) | `NOT_STARTED` | 0/3 | Pending | Multi-stage builds, non-root distroless, docker-compose |
| **19** | [Observability & Telemetry](./19-observability/README.md) | `NOT_STARTED` | 0/3 | Pending | Structured logging (`slog`), Prometheus metrics, tracing |
| **20** | [Microservice Architecture](./20-microservice/README.md) | `NOT_STARTED` | 0/3 | Pending | Domain-driven design, ports & adapters, clean internal layout |
| **21** | [Final Capstone Integration](./capstone/integration-service/README.md) | `NOT_STARTED` | 0/1 | Pending | Production integration service (Postgres, Redis, Workers) |

---

## 🏗️ Applied Projects Portfolio

In addition to modular stage exercises, this repository builds five standalone, production-styled applications:

* **[Project 1: `go-notes`](./projects/01-cli-tool/)** — Robust CLI note-taking tool with JSON persistence, flag parsing, and search.
* **[Project 2: `file-worker`](./projects/02-file-processor/)** — Concurrent multi-threaded file analyzer with bounded worker pools, rate limiting, and graceful cancellation.
* **[Project 3: `task-service`](./projects/03-rest-api/)** — Production REST API with PostgreSQL, connection pooling, migrations, and dockerized integration tests.
* **[Project 4: `event-dispatcher`](./projects/04-concurrent-worker/)** — High-throughput event pipeline with fan-out/fan-in workers, buffer management, and race protection.
* **[Project 5: `order-processor`](./projects/05-microservice/)** — Clean-architecture microservice featuring idempotency, transactions, structured logging, and health probes.
* **[Capstone: `integration-service`](./capstone/integration-service/)** — Full production integration engine connecting HTTP, PostgreSQL, Redis, background worker queues, and simulated external APIs.

---

## 🤖 The Specialist Agent System

Governing this repository is a specialized multi-agent architecture in [`.agents/`](./.agents/):

```text
.agents/
├── README.md                 # Agent architecture & interaction protocols
├── AGENTS.md                 # Antigravity IDE workspace rules
├── go-mentor.md              # Socratic guidance, progressive hints, no spoiler solutions
├── go-reviewer.md            # Idiomatic code reviewer ("What, Why, Idiomatic, Remember")
├── exercise-generator.md     # Hands-on exercises, test harnesses, separated solutions
├── project-architect.md      # Evolutionary architecture, domain design, ADR authoring
├── testing-engineer.md       # Table tests, race detection, benchmarks, fuzz testing
├── documentation-engineer.md # Memory-first cheat sheets, mental models, diagrams
├── git-github.md             # Conventional commits, CI workflows, portfolio integrity
├── devops-engineer.md        # Distroless Dockerfiles, Compose stacks, health checks
└── learning-tracker.md       # Evidence-based progress auditor in progress.yaml
```

### Mentoring Command Palette
Interact with the agent system using these command triggers:
- `start` — Get oriented and start the next milestone.
- `hint` — Receive progressive hints when stuck (without giving away code).
- `review` — Request an in-depth Go idiomatic review of your code.
- `solution` — Reveal the reference solution only after your attempt.
- `quiz` — Challenge your conceptual knowledge with targeted questions.
- `progress` — Audit repository tests and update `progress.yaml`.
- `next` — Advance to the next stage when all completion criteria are verified.

---

## 🧠 Memory-First Knowledge System

Located in [`learning/`](./learning/):
* **[Glossary](./learning/glossary.md)** — High-precision definitions with "Remember it like this" anchors.
* **[Mistakes Log](./learning/mistakes.md)** — Catalog of classic Go anti-patterns, traps, and debugging wisdom.
* **[Syntax Cheat Sheet](./learning/cheat-sheets/go-syntax.md)** — Rapid lookup for syntax, declarations, and control flow.
* **[Go Commands Cheat Sheet](./learning/cheat-sheets/go-commands.md)** — Essential Go CLI toolchain operations.

---

## 🚀 Quickstart & Verification

Verify your Go development environment:

```powershell
# Verify Go toolchain
go version

# Run all tests with race detector
go test -v -race ./...

# Verify code formatting
go fmt ./...

# Static analysis
go vet ./...
```
