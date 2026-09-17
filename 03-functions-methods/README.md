# Stage 03: Functions & Methods

## 🎯 Goal
Master Go's approach to procedural and object-oriented programming: first-class functions, lexical closures, variadic parameters, the precise lifecycle and LIFO mechanics of `defer`, and the critical distinction between value receivers and pointer receivers on custom types.

---

## 🧠 Core Concepts

### 1. First-Class Functions & Closures
In Go, functions are first-class citizens: they can be assigned to variables, passed as arguments to other functions, and returned from functions.
A **closure** is a function value that references variables outside its body. The closure binds to these variables; if the variable outlives the enclosing function, Go's escape analysis automatically moves the variable from the stack to the heap.

```go
func Counter() func() int {
    count := 0
    return func() int {
        count++ // captures 'count' by reference
        return count
    }
}
```

### 2. The `defer` Statement Mechanics
`defer` defers the execution of a function until the surrounding function returns.
- **LIFO Order**: Multiple deferred calls are pushed onto an internal stack and executed in Last-In, First-Out order.
- **Evaluation Timing**: Arguments to the deferred function are evaluated **immediately** when the `defer` statement is reached, not when the function actually executes!

```go
func example() {
    i := 0
    defer fmt.Println("deferred:", i) // Evaluates i=0 immediately!
    i = 10
    fmt.Println("current:", i)
}
// Output:
// current: 10
// deferred: 0
```

### 3. Value Receivers vs. Pointer Receivers
Go does not have classes, but you can define methods on custom types:
- **Value Receiver `(t T)`**: The method operates on a full copy of `T`. Any mutations made inside the method are discarded when the method finishes.
- **Pointer Receiver `(t *T)`**: The method operates on the actual memory address of `T`. Mutations alter the caller's struct, and pointer receivers avoid copying overhead for large data structures.

> **Rule of Thumb**:
> If any method of a type requires a pointer receiver, make **all** methods of that type use pointer receivers for consistency across the type's method set.

---

## 🔍 Mental Model

### Value Receiver vs. Pointer Receiver
```text
Value Receiver (t Type)          Pointer Receiver (t *Type)
┌───────────────────────┐        ┌───────────────────────┐
│ Original Struct in Mem│        │ Original Struct in Mem│
│       [ Name: "A" ]   │        │   [ Name: "A" ] <─────┐
└───────────────────────┘        └───────────────────────┤
            │ (photocopy)                                │ (gives address)
            ▼                                            │
┌───────────────────────┐                                │
│ Copy Inside Method    │                        Method reads/mutates
│       [ Name: "B" ]   │                        directly at address!
└───────────────────────┘
(Original stays "A")
```

---

## 💻 Examples
See executable code in:
- [`examples/closures_and_defer/main.go`](./examples/closures_and_defer/main.go)
- [`examples/methods_receivers/main.go`](./examples/methods_receivers/main.go)

---

## 🧪 Exercises

### 1. Functional Transformations Pipeline (Tier 1 - Easy)
- **Path**: [`exercises/01-function-pipelines/`](./exercises/01-function-pipelines/)
- **Concepts**: First-class functions, Map/Filter/Reduce combinators for integer slices.
- **Run tests**:
  ```powershell
  go test -v ./03-functions-methods/exercises/01-function-pipelines
  ```

### 2. Rate Limiter Closure (Tier 2 - Medium)
- **Path**: [`exercises/02-rate-limiter-closure/`](./exercises/02-rate-limiter-closure/)
- **Concepts**: State encapsulation with closures, time-window sliding, returning dynamic validator functions.
- **Run tests**:
  ```powershell
  go test -v ./03-functions-methods/exercises/02-rate-limiter-closure
  ```

### 3. Bank Account State Machine (Tier 3 - Hard)
- **Path**: [`exercises/03-bank-account-methods/`](./exercises/03-bank-account-methods/)
- **Concepts**: Custom types, value vs pointer receivers, immutable transaction ledgers, balance invariant enforcement.
- **Run tests**:
  ```powershell
  go test -v ./03-functions-methods/exercises/03-bank-account-methods
  ```

---

## 🧩 Challenge

Build a **Composable Middleware Pipeline** in [`challenges/middleware-pipeline/`](./challenges/middleware-pipeline/):
- Define `type HandlerFunc func(string) string` and `type Middleware func(HandlerFunc) HandlerFunc`.
- Build middlewares for logging execution time, input sanitization, and authorization token checking.
- Implement `Chain(middlewares ...Middleware) Middleware` to wrap handlers cleanly.

---

## ⚠️ Common Mistakes

1. **Defer in Long Loops**:
   ```go
   // ANTI-PATTERN: File handles stay open until entire function exits!
   for _, filename := range files {
       f, _ := os.Open(filename)
       defer f.Close() // Memory / file descriptor exhaustion!
   }
   // FIX: Wrap inside a helper function, or close explicitly without defer
   ```
2. **Mutating on Value Receiver**:
   ```go
   type Counter struct { count int }
   func (c Counter) Increment() { c.count++ } // Mutates a local copy!
   ```

---

## 🔄 Review

1. When does a deferred function actually execute?
2. In what order do multiple deferred calls execute?
3. When are arguments to a deferred function evaluated?
4. When should a method use a pointer receiver instead of a value receiver?
5. How does Go determine whether a variable captured by a closure lives on the stack or escapes to the heap?

---

## ✅ Completion Criteria

- [ ] All 3 exercises pass unit tests (`go test -v ./...`).
- [ ] Challenge in `challenges/middleware-pipeline/` passes tests.
- [ ] Code passes `go vet` and `go fmt`.
- [ ] Reviewed by **Go Reviewer**.

---

## 🔗 Connections
- **Prerequisites**: [Stage 02: Data Structures](../02-data-structures/README.md).
- **Next Stage**: [Stage 04: Structs & Interfaces](../04-structs-interfaces/README.md) (Implicit interface satisfaction, polymorphism).
