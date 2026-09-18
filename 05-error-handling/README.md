# Stage 05: Error Handling

## 🎯 Goal
Master Go's core philosophy that **errors are normal values, not exceptional control flow**: the built-in `error` interface, sentinel error constants, contextual error wrapping using `%w`, deep chain inspection with `errors.Is` and `errors.As`, building domain-specific custom error types, and understanding when to use `panic`/`recover`.

---

## 🧠 Core Concepts

### 1. Errors as Values
In Go, there are no exceptions (`try / catch / throw`). An error is an ordinary value implementing the minimal built-in interface:
```go
type error interface {
    Error() string
}
```
Functions return errors as the final return value. Callers explicitly check `if err != nil` and either handle the error locally or wrap and propagate it up the stack.

### 2. Sentinel Errors & `errors.Is()`
A **sentinel error** is a package-level variable representing a known, specific failure condition (e.g. `io.EOF`, `sql.ErrNoRows`, `os.ErrNotExist`).
To check if an error matches a sentinel error—even after it has been wrapped—always use `errors.Is(err, target)` instead of `==`:
```go
var ErrNotFound = errors.New("record not found")

// Later:
if errors.Is(err, ErrNotFound) {
    // Matches even if err was wrapped: fmt.Errorf("fetching user 42: %w", ErrNotFound)
}
```

### 3. Contextual Error Wrapping (`%w`)
When propagating errors, prepend contextual information about what the current function was trying to do:
```go
data, err := os.ReadFile(filename)
if err != nil {
    return fmt.Errorf("reading configuration from %s: %w", filename, err)
}
```
> **Notice**:
> Use `%w` to wrap the error so callers can unpack it with `errors.Is` or `errors.As`.
> Use `%v` if you deliberately want to obscure or hide internal error details from the caller.

### 4. Custom Error Structs & `errors.As()`
When an error requires structured metadata (such as an HTTP status code, database error code, or offending field name), define a custom struct:
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed on '%s': %s", e.Field, e.Message)
}
```
To extract a specific custom error type from a wrapped chain, use `errors.As()`:
```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Printf("Offending field was: %s\n", valErr.Field)
}
```

---

## 🔍 Mental Model

Think of error wrapping as nesting **Russian Matryoshka Dolls**:

```text
┌─────────────────────────────────────────────────────────────┐
│ Outer Layer: "processing checkout order 1042"               │
│   ┌───────────────────────────────────────────────────────┐ │
│   │ Middle Layer: "charging payment with stripe"          │ │
│   │   ┌─────────────────────────────────────────────────┐ │ │
│   │   │ Root Cause: ErrNetworkTimeout ("dial tcp timeout")│ │ │
│   │   └─────────────────────────────────────────────────┘ │ │
│   └───────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```
- `errors.Is(err, ErrNetworkTimeout)` opens all dolls to see if the core doll is `ErrNetworkTimeout`.
- `errors.As(err, &target)` opens all dolls to see if any doll is of type `*StripeError`.

---

## 💻 Examples
See executable code in:
- [`examples/sentinel_and_custom_errors/main.go`](./examples/sentinel_and_custom_errors/main.go)
- [`examples/wrapping_is_as/main.go`](./examples/wrapping_is_as/main.go)

---

## 🧪 Exercises

### 1. Sentinel Errors & Query Engine (Tier 1 - Easy)
- **Path**: [`exercises/01-sentinel-errors/`](./exercises/01-sentinel-errors/)
- **Concepts**: Defining package-level sentinel errors (`ErrUserNotFound`, `ErrUnauthorized`, `ErrDuplicateEmail`), error wrapping with `%w`, and testing with `errors.Is`.
- **Run tests**:
  ```powershell
  go test -v ./05-error-handling/exercises/01-sentinel-errors
  ```

### 2. Structured Domain Validation Errors (Tier 2 - Medium)
- **Path**: [`exercises/02-custom-error-types/`](./exercises/02-custom-error-types/)
- **Concepts**: Custom struct error types, implementing `Error() string`, extracting fields with `errors.As()`.
- **Run tests**:
  ```powershell
  go test -v ./05-error-handling/exercises/02-custom-error-types
  ```

### 3. Resilient Retry with Error Classification (Tier 3 - Hard)
- **Path**: [`exercises/03-resilient-retry/`](./exercises/03-resilient-retry/)
- **Concepts**: Classifying errors into transient (retryable) vs permanent (non-retryable), exponential backoff, maximum retry attempts.
- **Run tests**:
  ```powershell
  go test -v ./05-error-handling/exercises/03-resilient-retry
  ```

---

## 🧩 Challenge

Build an **Aggregate Multi-Error Collector** in [`challenges/multierror-collector/`](./challenges/multierror-collector/):
- Implement `type MultiError struct` satisfying `error` and the Go 1.20+ `Unwrap() []error` interface.
- Collect multiple independent validation or batch errors without failing on the first error.
- Format clean multi-line error summaries.

---

## ⚠️ Common Mistakes

1. **Comparing Wrapped Errors with `==`**:
   ```go
   // ANTI-PATTERN: Fails if err was wrapped with fmt.Errorf("%w")!
   if err == ErrNotFound { ... }
   
   // IDIOMATIC:
   if errors.Is(err, ErrNotFound) { ... }
   ```
2. **Ignoring Errors with the Blank Identifier**:
   ```go
   _ = json.Unmarshal(data, &v) // ANTI-PATTERN: Silent failures!
   ```
3. **Using Panic for Regular Error Handling**: In Go, `panic` is only for unrecoverable programmer errors (e.g. nil pointer, corrupted state). Network disconnects, bad user input, and missing files are ordinary errors.

---

## 🔄 Review

1. What is the difference between `%w` and `%v` in `fmt.Errorf()`?
2. When should you use `errors.Is()` vs `errors.As()`?
3. What is a sentinel error, and why should it be declared as an unexported or exported package variable?
4. What happens if a function panics and no `recover()` is called?

---

## ✅ Completion Criteria

- [ ] All 3 exercises pass unit tests (`go test -v ./...`).
- [ ] Challenge in `challenges/multierror-collector/` passes tests.
- [ ] Code passes `go vet` and `go fmt`.
- [ ] Progress updated in `learning/progress.yaml`.

---

## 🔗 Connections
- **Prerequisites**: [Stage 04: Structs & Interfaces](../04-structs-interfaces/README.md).
- **Next Stage**: [Stage 06: Packages & Modules](../06-packages-modules/README.md) (`internal/` boundaries, module versioning).
