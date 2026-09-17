# Stage 04: Structs & Interfaces

## 🎯 Goal
Master Go's composition-over-inheritance paradigm: struct memory layout, struct embedding, implicit interface satisfaction, small composable interfaces, the internal `(type, value)` interface tuple, and how to avoid the infamous *"typed nil interface"* trap.

---

## 🧠 Core Concepts

### 1. Structs, Memory Alignment & Padding
Structs in Go are contiguous blocks of memory. The order of fields matters due to CPU word alignment:
```go
// 24 bytes due to padding:
type Unoptimized struct {
    a bool    // 1 byte + 7 bytes padding
    b int64   // 8 bytes
    c bool    // 1 byte + 7 bytes padding
}

// 16 bytes:
type Optimized struct {
    b int64   // 8 bytes
    a bool    // 1 byte
    c bool    // 1 byte + 6 bytes padding
}
```

### 2. Composition Over Inheritance (Embedding)
Go does not have class hierarchies or `extends`. Instead, types achieve reuse via **embedding**:
- An embedded type's methods and fields are "promoted" to the outer type.
- The outer type can override promoted methods without breaking the inner type.

### 3. Implicit Interface Satisfaction
An interface defines a set of method signatures. A concrete type implements an interface simply by implementing those methods:
- No `implements` keyword.
- Allows retroactive interface definition: consumers can define interfaces that third-party types already satisfy!
- **Go Proverb**: *"Accept interfaces, return structs."*

### 4. Interface Internals & The "Typed Nil" Trap
An interface variable is a 2-word data structure:
```text
Interface Value in Memory
┌─────────────────────────┬─────────────────────────┐
│     Type Descriptor     │      Value Pointer      │
│      (concrete type)    │    (address of data)    │
└─────────────────────────┴─────────────────────────┘
```
An interface is `nil` **only** if both the Type Descriptor and the Value Pointer are `nil`!
If you assign a `nil` pointer of a concrete type to an interface, the interface's Type Descriptor is non-nil, so `iface == nil` evaluates to **`false`**!

```go
// THE TYPED NIL TRAP:
var myErr *CustomError = nil // concrete pointer is nil
var err error = myErr        // interface gets type *CustomError
if err != nil {
    // THIS BRANCH EXECUTES! err is not nil because its type slot is populated!
}
```

---

## 🔍 Mental Model

Think of an interface as a two-pin adapter:

```text
┌──────────────────────────────────────────────┐
│                Interface Box                 │
│                                              │
│  [ Type: *PostgresDB ]    [ Ptr: 0x1400a0 ]  │
└───────────────────────┬──────────────────────┘
                        │
                        ▼
                Actual Instance in Heap
```
If you plug in a cord that points to nothing (`nil`), the adapter still holds the plug specification. It is not an empty socket.

---

## 💻 Examples
See executable code in:
- [`examples/struct_embedding/main.go`](./examples/struct_embedding/main.go)
- [`examples/interfaces_polymorphism/main.go`](./examples/interfaces_polymorphism/main.go)

---

## 🧪 Exercises

### 1. Geometric Shapes & Polymorphism (Tier 1 - Easy)
- **Path**: [`exercises/01-shapes-geometry/`](./exercises/01-shapes-geometry/)
- **Concepts**: `Shape` interface (`Area()`, `Perimeter()`), `Circle`, `Rectangle`, total area calculation.
- **Run tests**:
  ```powershell
  go test -v ./04-structs-interfaces/exercises/01-shapes-geometry
  ```

### 2. Multi-Provider Payment Gateway (Tier 2 - Medium)
- **Path**: [`exercises/02-payment-gateway/`](./exercises/02-payment-gateway/)
- **Concepts**: `PaymentProcessor` interface, concrete implementations (`CreditCard`, `PayPal`, `Crypto`), unified transaction dispatcher.
- **Run tests**:
  ```powershell
  go test -v ./04-structs-interfaces/exercises/02-payment-gateway
  ```

### 3. Dynamic Plugin Lifecycle Engine (Tier 3 - Hard)
- **Path**: [`exercises/03-plugin-pipeline/`](./exercises/03-plugin-pipeline/)
- **Concepts**: Interface segregation, optional interface capability checking using type assertions (`p, ok := plugin.(HealthChecker)`).
- **Run tests**:
  ```powershell
  go test -v ./04-structs-interfaces/exercises/03-plugin-pipeline
  ```

---

## 🧩 Challenge

Build an extensible **Multi-Channel Notification Engine** in [`challenges/notifying-service/`](./challenges/notifying-service/):
- Define `Notifier` interface.
- Implement `EmailNotifier`, `SMSNotifier`, and `WebhookNotifier`.
- Build a `MultiNotifier` composite aggregator with retry policies and audit mock records.

---

## ⚠️ Common Mistakes

1. **The Typed Nil Error Return**:
   ```go
   func DoWork() error {
       var err *MyCustomError = nil
       return err // Anti-pattern: returns a non-nil error interface!
   }
   // FIX: Return explicit nil
   func DoWork() error {
       return nil
   }
   ```
2. **Premature Interface Abstraction**: Creating an interface for every single struct before having 2+ implementations. In Go, define interfaces where they are consumed, not alongside implementations.

---

## 🔄 Review

1. Under what condition is an interface variable considered `nil` in Go?
2. What is the difference between struct embedding and traditional OOP inheritance?
3. Why does the Go community prefer small interfaces (e.g. `io.Reader` with 1 method) over fat interfaces?
4. How do you test whether an interface value implements an optional interface without panicking?

---

## ✅ Completion Criteria

- [ ] All 3 exercises pass unit tests (`go test -v ./...`).
- [ ] Challenge in `challenges/notifying-service/` passes tests.
- [ ] You can explain why `var e *MyErr = nil; var err error = e; err != nil` evaluates to true.
- [ ] Progress updated in `learning/progress.yaml`.

---

## 🔗 Connections
- **Prerequisites**: [Stage 03: Functions & Methods](../03-functions-methods/README.md).
- **Next Stage**: [Stage 05: Error Handling](../05-error-handling/README.md) (Custom errors, wrapping, `errors.Is`, `errors.As`).
