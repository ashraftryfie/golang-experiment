# Stage 06: Packages & Modules

## 🎯 Goal
Master Go's module and package architecture: visibility rules (exported vs unexported), package naming conventions, compiler-enforced `internal/` access boundaries, `init()` function execution order, avoiding and resolving cyclic dependencies, and `go.mod`/`go.sum` dependency mechanics.

---

## 🧠 Core Concepts

### 1. Package Naming & Organization
- Every `.go` source file starts with `package <name>`.
- All files in the same directory **must** belong to the same package (with the single exception of `package <name>_test` for black-box testing).
- **Go Convention**: Package names should be short, lowercase, single-word nouns without underscores or camelCase (e.g. `user`, `json`, `http`, not `user_service` or `userManagement`).
- Avoid generic catch-all packages like `util`, `common`, or `helpers`. Package by domain or capability.

### 2. Visibility & Export Rules
Visibility in Go is governed exclusively by the first letter of an identifier:
- **Exported (Public)**: Starts with an uppercase letter (`User`, `NewService`, `ErrNotFound`). Accessible outside the package.
- **Unexported (Private)**: Starts with a lowercase letter (`user`, `validateInput`, `secretKey`). Scoped strictly to the package.

### 3. The Special `internal/` Boundary
Go's compiler natively enforces encapsulation through any directory named `internal/`:
- Code inside `.../internal/foo` can **only** be imported by packages within the parent tree rooted at the directory containing `internal`.
- External modules or sibling packages outside that parent directory cannot import it, preventing external code from coupling to your private implementation details.

```text
my-service/
├── cmd/
│   └── api/ (can import internal/db)
├── internal/
│   └── db/  (PRIVATE to my-service; external projects cannot import this!)
└── pkg/
    └── mathutil/ (PUBLIC to the world)
```

### 4. Initialization & `init()` Functions
1. First, package-level variables are evaluated in dependency order.
2. Next, any `init()` functions are executed. A package can have multiple `init()` functions across multiple files.
3. `init()` takes no arguments and returns no values.
4. **Anti-Pattern Warning**: Avoid heavy logic, network calls, or non-deterministic state in `init()`. Prefer explicit constructors (`New()`).

### 5. Cyclic Imports & Decoupling
Go strictly forbids cyclic dependencies (`import cycle not allowed`):
```text
Package A imports Package B
    ▲                 │
    │                 ▼
    └── Package B imports Package A  ❌ (COMPILE ERROR)
```
**How to Break Cycles**:
1. **Interface Inversion**: Define the interface where the service consumes it, not where it is implemented.
2. **Move Shared Types**: Extract the shared domain models or types into a lower-level leaf package (e.g. `domain` or `types`).
3. **Dependency Injection**: Pass the dependency into the constructor function.

---

## 🔍 Mental Model

Think of package dependencies as a **Directed Acyclic Graph (DAG)** of building blocks:

```text
           cmd/server (Binary Entrypoint)
               │                 │
               ▼                 ▼
         internal/api      internal/worker
               │                 │
               └────────┬────────┘
                        ▼
                 internal/domain (Pure Business Logic)
                        │
                        ▼
                 standard library (Leaf Nodes)
```
Dependencies flow downward. Lower layers never know about upper layers.

---

## 💻 Examples
See executable code in:
- [`examples/visibility_and_internal/`](./examples/visibility_and_internal/)
- [`examples/init_execution_order/main.go`](./examples/init_execution_order/main.go)

---

## 🧪 Exercises

### 1. API Encapsulation & Export Auditing (Tier 1 - Easy)
- **Path**: [`exercises/01-package-export/`](./exercises/01-package-export/)
- **Concepts**: Refactoring leaked internal state into unexported fields with exported getters and constructor functions (`NewConfig()`).
- **Run tests**:
  ```powershell
  go test -v ./06-packages-modules/exercises/01-package-export
  ```

### 2. Internal Security Boundary (Tier 2 - Medium)
- **Path**: [`exercises/02-internal-boundary/`](./exercises/02-internal-boundary/)
- **Concepts**: Structuring a token service that uses an unexported `internal/hasher` package, exposing clean public signing functions.
- **Run tests**:
  ```powershell
  go test -v ./06-packages-modules/exercises/02-internal-boundary/...
  ```

### 3. Decoupling Cyclic Imports (Tier 3 - Hard)
- **Path**: [`exercises/03-cyclic-decoupling/`](./exercises/03-cyclic-decoupling/)
- **Concepts**: Breaking simulated cyclic dependencies between `users` and `orders` using interface inversion and constructor injection.
- **Run tests**:
  ```powershell
  go test -v ./06-packages-modules/exercises/03-cyclic-decoupling/...
  ```

---

## 🧩 Challenge

Build a **Modular Architecture Scaffold** in [`challenges/modular-monolith/`](./challenges/modular-monolith/):
- Construct a layered structure: `internal/domain` (types), `internal/store` (storage interface + memory implementation), and `pkg/validator` (reusable helpers).
- Demonstrate clean cross-package wiring in `main.go`.

---

## ⚠️ Common Mistakes

1. **Naming stutter**:
   ```go
   // ANTI-PATTERN:
   package user
   type UserUser struct { ... } // Caller writes: user.UserUser
   
   // IDIOMATIC:
   package user
   type User struct { ... }     // Caller writes: user.User
   ```
2. **Polluting Global State in `init()`**: Initializing database connections or HTTP clients inside `init()` prevents proper error handling, testing, and configuration passing. Always use explicit constructors.
3. **Circular Dependencies**: Trying to make two packages depend on each other. When this happens, your package boundaries are wrong.

---

## 🔄 Review

1. How does Go determine if a function, struct, or variable is accessible outside its package?
2. What special compiler rule applies to any directory named `internal/`?
3. Why does Go forbid circular package imports?
4. When does the `init()` function execute relative to `main()`?

---

## ✅ Completion Criteria

- [ ] All 3 exercises pass unit tests (`go test -v ./...`).
- [ ] Challenge in `challenges/modular-monolith/` passes tests.
- [ ] Code passes `go vet` and `go fmt`.
- [ ] Progress updated in `learning/progress.yaml`.

---

## 🔗 Connections
- **Prerequisites**: [Stage 05: Error Handling](../05-error-handling/README.md).
- **Next Stage**: [Stage 07: Standard Library Deep Dive](../07-standard-library/README.md) (`io.Reader`, `io.Writer`, streaming pipelines).
