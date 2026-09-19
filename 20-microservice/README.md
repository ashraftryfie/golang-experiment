# Stage 20: Microservice Architecture

Welcome to **Stage 20** of your Go mastery journey. Building production-grade microservices in Go demands more than writing HTTP handlers. It requires resilient software architecture that decouples business logic from external frameworks, databases, and third-party transport layers.

In this stage, you master **Clean Architecture (Hexagonal / Ports & Adapters)** and **Domain-Driven Design (DDD)** tailored specifically for Go idioms.

---

## 1. Hexagonal Architecture (Ports & Adapters)

The core principle of Hexagonal Architecture is **Dependency Inversion**:
> *High-level business rules must not depend on low-level technical mechanisms (databases, HTTP frameworks, third-party SDKs). Both must depend on abstractions (interfaces).*

```text
┌────────────────────────────────────────────────────────┐
│                        ADAPTERS                        │
│   HTTP Handlers, gRPC Services, CLI, Message Consumers │
│                           │                            │
│                           ▼                            │
│ ┌────────────────────────────────────────────────────┐ │
│ │                   SERVICE LAYER                    │ │
│ │             (Use-Case Orchestration)               │ │
│ │                         │                          │ │
│ │                         ▼                          │ │
│ │ ┌────────────────────────────────────────────────┐ │ │
│ │ │                 DOMAIN MODEL                   │ │ │
│ │ │    Entities, Value Objects, Domain Errors      │ │ │
│ │ └────────────────────────────────────────────────┘ │ │
│ │                         ▲                          │ │
│ │                         │                          │ │
│ │                     PORTS (Interfaces)             │ │
│ └─────────────────────────┬──────────────────────────┘ │
│                           │                            │
│                           ▼                            │
│                        ADAPTERS                        │
│     PostgreSQL, Redis Cache, Kafka Producer, Stripe    │
└────────────────────────────────────────────────────────┘
```

### The Go Axiom: "Accept Interfaces, Return Structs"
In Go, interfaces are satisfied **implicitly**. You do not declare `implements`. Therefore:
1. **Define interfaces at the point of consumption** (where they are used, not where they are implemented).
2. The domain and service layers define the **Ports** (`Repository`, `Notifier`, `PaymentGateway`).
3. The persistence, transport, and external SDK packages provide the **Adapters** (`PostgresRepo`, `HTTPHandler`, `StripeClient`).

---

## 2. Domain-Driven Design (DDD) Fundamentals

### A. Entities vs. Value Objects
* **Entity**: An object defined by its identity, not its attributes. It has a lifecycle and mutates over time (e.g., `Order`, `Account`, `User`).
* **Value Object**: An immutable object defined entirely by its properties. Two value objects with the same values are identical (e.g., `Money{Amount: 100, Currency: "USD"}`).

```go
// Value Object: Immutable representation of money
type Money struct {
    amount   int64 // Stored in minor currency units (cents) to avoid floating point drift
    currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
    if amount < 0 {
        return Money{}, ErrNegativeAmount
    }
    if currency == "" {
        return Money{}, ErrInvalidCurrency
    }
    return Money{amount: amount, currency: currency}, nil
}
```

### B. Encapsulation & Domain Invariants
Never export struct fields if external code can bypass validation. Expose behaviors and constructor functions that enforce business invariants:

```go
type Order struct {
    id        string
    status    OrderStatus
    items     []LineItem
    total     Money
}

func (o *Order) Cancel() error {
    if o.status == OrderStatusShipped {
        return ErrCannotCancelShippedOrder
    }
    o.status = OrderStatusCancelled
    return nil
}
```

---

## 3. Standard Internal Layout for Go Microservices

Go standard layout separates concerns within `internal/` to prevent external packages from importing internal implementations:

```text
20-microservice/
├── README.md
├── examples/
│   └── hexagonal_wiring/main.go  # End-to-end DI and orchestration demo
├── exercises/
│   ├── 01-domain-entities        # Value objects, entities, and domain invariants
│   ├── 02-ports-and-adapters     # Consumer-driven interfaces & adapter implementations
│   └── 03-service-layer-di       # Use-case orchestration & constructor injection
├── solutions/
│   ├── 01-domain-entities
│   ├── 02-ports-and-adapters
│   └── 03-service-layer-di
└── challenges/
    └── order-service-hexagonal   # Production-grade hexagonal microservice
```

---

## 4. Key Rules for Idiomatic Architecture

1. **Zero External Dependencies in Domain**: `domain` must only import standard library packages. Never import SQL drivers, ORMs, or HTTP libraries into the domain core.
2. **Constructor-Based Dependency Injection**: Use clear constructor functions (`NewService(repo Repository, pub Publisher) *Service`) rather than reflection-heavy DI frameworks.
3. **Context as First Parameter**: Every port method performing I/O must accept `ctx context.Context` as its first parameter.
