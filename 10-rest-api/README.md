# Stage 10: REST API Engineering

Welcome to **Stage 10** of your Go engineering journey. In this stage, you bridge raw HTTP handlers into production-grade REST APIs: resource modeling, idempotency semantics, robust input validation, thread-safe in-memory storage, and RFC 7807 standardized error formatting.

---

## 🧭 Mental Model

Building REST APIs in Go requires adhering to clear semantic contracts:

### 1. HTTP Verbs & Idempotency Semantics
| Method | Idempotent | Safe | Expected Status Codes | Typical Usage |
| :--- | :---: | :---: | :--- | :--- |
| `GET` | Yes | Yes | `200 OK`, `404 Not Found` | Fetch resource or list |
| `POST` | **No** | No | `201 Created` (with `Location`), `400`, `422` | Create resource |
| `PUT` | Yes | No | `200 OK`, `204 No Content` | Complete replacement |
| `PATCH`| **No** | No | `200 OK` | Partial update |
| `DELETE`| Yes | No | `204 No Content`, `404 Not Found` | Delete resource |

### 2. Robust JSON Request Decoding
Never use raw unmarshaling on untrusted HTTP request bodies. Always enforce:
```go
dec := json.NewDecoder(r.Body)
dec.DisallowUnknownFields() // Reject unexpected client fields

if err := dec.Decode(&payload); err != nil {
    // Distinguish syntax errors from schema mismatches
}
```

### 3. RFC 7807: Standardized Problem Details
Instead of ad-hoc error formats, production microservices use the IETF standard **RFC 7807** (`application/problem+json`):
```json
{
  "type": "https://api.example.com/errors/invalid-parameters",
  "title": "Invalid Request Parameters",
  "status": 422,
  "detail": "The 'price' field must be greater than zero.",
  "instance": "/items/create"
}
```

### 4. Thread-Safe In-Memory Stores
Because every HTTP request in Go is processed in its own independent goroutine, shared state (like an in-memory repository) MUST be protected using concurrency primitives like `sync.RWMutex`:
```go
type ItemStore struct {
    mu    sync.RWMutex
    items map[string]Item
}
```

---

## 🎯 Learning Objectives

By the end of this stage, you will:
- [x] Design idiomatic REST resource endpoints according to HTTP standards.
- [x] Safely decode JSON request bodies with `DisallowUnknownFields()`.
- [x] Validate domain business rules and return descriptive `422 Unprocessable Entity` or `400 Bad Request`.
- [x] Implement thread-safe CRUD operations using `sync.RWMutex`.
- [x] Structure standardized error responses following RFC 7807.
- [x] Parse query parameters for pagination (`?page=1&limit=20`) and filters.

---

## 🗂️ Stage Structure

```
10-rest-api/
├── README.md                           # Stage curriculum & architecture
├── examples/
│   ├── validation/main.go              # Request validation & DisallowUnknownFields
│   └── crud_service/main.go            # Thread-safe in-memory REST service
├── exercises/
│   ├── 01-rest-models/                 # Domain model validation & error mapping
│   ├── 02-crud-handlers/               # Thread-safe REST CRUD HTTP handlers
│   └── 03-rfc7807-errors/              # RFC 7807 Problem Details response builder
├── solutions/
│   ├── 01-rest-models/                 # Reference implementation
│   ├── 02-crud-handlers/               # Reference implementation
│   └── 03-rfc7807-errors/              # Reference implementation
└── challenges/
    └── task-manager-api/               # Complete Task REST API with filtering & pagination
```

---

## 🧪 Verification Commands

```powershell
# Run all tests in this stage
go test -v ./10-rest-api/...

# Vet and check format
go vet ./10-rest-api/...
go fmt ./10-rest-api/...
```

---

## 🔗 Connections
- **Prerequisites**: [Stage 09: HTTP Fundamentals](../09-http/README.md).
- **Next Stage**: [Stage 11: PostgreSQL & Storage](../11-postgresql/README.md) (`database/sql`, connection pooling, transactions).
