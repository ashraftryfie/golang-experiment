# Project 5: `order-processor` (Order Processing Microservice)

A clean-architecture microservice featuring idempotency keys, state machine validations, structured logging, and health probes.

## Features
- **Idempotency Safeguard**: Every transaction accepts an `Idempotency-Key` header/parameter to eliminate duplicate execution during network retries.
- **State Machine Transitions**: Guards order statuses (`Pending` -> `Confirmed` -> `Completed` / `Cancelled`).
- **Telemetry & Health Probes**: Exposes readiness and liveness checks for Kubernetes deployment.

## Usage
```powershell
go test -v ./projects/05-microservice/...
```
