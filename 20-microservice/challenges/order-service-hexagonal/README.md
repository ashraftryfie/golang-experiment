# Challenge: Hexagonal Order Microservice

## Objective
Design and implement a complete production-grade microservice component following **Hexagonal Architecture (Ports & Adapters)** and **Domain-Driven Design (DDD)** in Go.

## Architecture

```text
┌─────────────────────────────────────────────────────────────┐
│                       HTTP ADAPTER                          │
│        POST /orders, POST /orders/{id}/pay, GET /orders/{id}│
│                             │                               │
│                             ▼                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │                    ORDER SERVICE                      │  │
│  │    (Orchestration, Invariants, Payments, Events)      │  │
│  │                          │                            │  │
│  │                          ▼                            │  │
│  │  ┌─────────────────────────────────────────────────┐  │  │
│  │  │                 DOMAIN MODEL                    │  │  │
│  │  │  Order, OrderItem, OrderStatus, State Machine   │  │  │
│  │  └─────────────────────────────────────────────────┘  │  │
│  │                          ▲                            │  │
│  │                          │                            │  │
│  │       PORTS: OrderRepository, PaymentGateway,         │  │
│  │              EventDispatcher                          │  │
│  └──────────────────────────┬────────────────────────────┘  │
│                             │                               │
│                             ▼                               │
│      ADAPTERS: InMemoryRepo, MockPayment, MockDispatcher    │
└─────────────────────────────────────────────────────────────┘
```

## Requirements
1. **Domain Model**:
   - `OrderStatus`: `StatusPending`, `StatusPaid`, `StatusShipped`, `StatusCancelled`.
   - `Order` state transitions:
     - `Pending` -> `Paid` or `Cancelled`.
     - `Paid` -> `Shipped` or `Cancelled`.
     - `Shipped` cannot be cancelled (`ErrCannotCancelShipped`).
     - Cannot transition to `Paid` if already `Paid` (`ErrAlreadyPaid`).
   - `TotalCents()` computed dynamically from items.

2. **Ports**:
   - `OrderRepository`: Thread-safe persistence port.
   - `PaymentGateway`: Payment processor port.
   - `EventDispatcher`: Outbound domain event publisher.

3. **Service Layer**:
   - `CreateOrder`: Validates items, initializes status to `StatusPending`, persists order, dispatches `order.created` event.
   - `PayOrder`: Fetches order, charges customer via `PaymentGateway`, transitions order to `StatusPaid`, persists, dispatches `order.paid`.
   - `ShipOrder`: Transitions order to `StatusShipped`, persists, dispatches `order.shipped`.
   - `CancelOrder`: Enforces cancellation invariants, transitions order to `StatusCancelled`, persists, dispatches `order.cancelled`.

4. **HTTP Adapter**:
   - Handlers parsing JSON requests, invoking `OrderService`, and returning clean JSON responses with proper HTTP status codes (`201 Created`, `200 OK`, `400 Bad Request`, `404 Not Found`).
