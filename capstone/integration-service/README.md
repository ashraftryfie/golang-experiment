# Capstone Project: Production Integration Service (`integration-service`)

Welcome to the **Final Capstone Integration Service** of the Go mastery path. This project brings together everything mastered across all 20 previous stages:
- **Clean Hexagonal Architecture**: Strict decoupling of domain models from transport and persistence adapters.
- **Concurrency & Worker Pools**: Bounded worker routines processing jobs asynchronously via channel brokers with rate-limiting and backoff.
- **Context & Resiliency**: Context-driven timeout propagation, graceful cancellation, and compensating transactions.
- **Observability**: Structured logging (`slog`), Prometheus-compatible metrics exposition (`/metrics`), and health check probes (`/healthz`).
- **Graceful Shutdown**: Zero-downtime shutdown capturing `SIGINT`/`SIGTERM`, draining in-flight worker jobs and terminating HTTP listeners cleanly.

---

## 🏗️ System Architecture

```text
                                  CLIENT
                                    │
                         ┌──────────┴──────────┐
                         │  HTTP REST Gateway  │
                         └──────────┬──────────┘
                                    │
               ┌────────────────────▼────────────────────┐
               │         INTEGRATION SERVICE             │
               │   (Use-Case Orchestration & Logic)      │
               └──────────┬───────────────────┬──────────┘
                          │                   │
                          ▼                   ▼
                 ┌────────────────┐   ┌───────────────┐
                 │  Job Storage   │   │  Queue Broker │
                 │     (Port)     │   │     (Port)    │
                 └────────┬───────┘   └───────┬───────┘
                          │                   │
                          │                   ▼
                          │         ┌───────────────────┐
                          │         │ Concurrent Worker │
                          │         │       Pool        │
                          │         └─────────┬─────────┘
                          │                   │
                          ▼                   ▼
                 ┌────────────────┐   ┌───────────────┐
                 │ InMemoryStore  │   │  External API │
                 │   (Adapter)    │   │  (Sync / Port)│
                 └────────────────┘   └───────────────┘
```

---

## 📡 API Specifications

| Method | Route | Description | Expected Status |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/jobs` | Enqueue a new asynchronous integration job | `201 Created` |
| `GET` | `/api/v1/jobs/{id}` | Fetch current job status, execution attempts, and payload | `200 OK` / `404 Not Found` |
| `GET` | `/healthz` | Health and readiness probe | `200 OK` |
| `GET` | `/metrics` | Prometheus metrics exposition endpoint | `200 OK` |

---

## ⚙️ Running and Testing

Run integration tests:
```powershell
go test -v ./capstone/integration-service/test/...
```

Run server standalone:
```powershell
go run ./capstone/integration-service/cmd/server
```
