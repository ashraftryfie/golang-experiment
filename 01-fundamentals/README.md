# Stage 01: Go Fundamentals

## 🎯 Goal
Master Go's static type system, zero-value guarantees, explicit conversions, control flow constructs (`if` with initializer, `switch`, `for`), functions with multiple return values, and write your first tested Go code.

---

## 🧠 Core Concepts
- **Static Typing & Type Safety**: Go is strictly typed. There is no implicit type coercion (you cannot implicitly add an `int` to a `float64`).
- **Zero Value Guarantee**: Declaring a variable without initialization assigns its type-safe zero value (`0`, `0.0`, `""`, `false`, `nil`).
- **Short Variable Declaration (`:=`)**: Used exclusively within function bodies for local variable declaration and type inference.
- **Constants (Typed vs Untyped)**: Untyped constants have arbitrary precision and adapt to the type of expressions they interact with until assigned to a typed variable.
- **Control Flow Minimalism**: Go has only one looping keyword: `for`.
- **Multiple Return Values**: Functions commonly return a result alongside an `error` value.

---

## 🔍 Mental Model

Think of Go's type system as physical containers of exact shapes:

```text
┌──────────────┐     ┌──────────────┐
│   int (64)   │ ≠  │  float64     │
│   [ 42 ]     │     │  [ 42.0 ]    │
└──────────────┘     └──────────────┘
       ▲                    ▲
       │                    │
       └── No auto-mixing ──┘
Must explicitly convert: float64(42)
```

In JavaScript or Python, numbers can silently change representation or produce unexpected coercions (`"5" + 3 = "53"`).
In Go, if types do not match exactly, the compiler stops immediately. This eliminates an entire class of runtime bugs.

---

## 💻 Examples

### Types, Zero Values & Conversions
See [`examples/types_and_vars/main.go`](./examples/types_and_vars/main.go):
```go
var count int            // Guaranteed 0
var temperature float64  // Guaranteed 0.0
var city string          // Guaranteed ""

// Explicit conversion is mandatory
total := float64(count) + temperature
```

### Control Flow
See [`examples/control_flow/main.go`](./examples/control_flow/main.go):
```go
// If with initializer (limits scope of 'result')
if result, err := divide(10, 2); err == nil {
    fmt.Println("Success:", result)
}

// Switch without expression (clean if-else alternative)
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "C"
}
```

---

## 🧪 Exercises

Exercises must be completed by you in the `exercises/` directory. Each exercise contains starter code with `TODO` markers and a pre-configured test suite.

### 1. Temperature Converter (Easy)
- **Path**: [`exercises/01-temperature-converter/`](./exercises/01-temperature-converter/)
- **Task**: Implement conversions between Celsius, Fahrenheit, and Kelvin with boundary validations.
- **Run tests**:
  ```powershell
  go test -v ./01-fundamentals/exercises/01-temperature-converter
  ```

### 2. FizzBuzz Variations (Medium)
- **Path**: [`exercises/02-fizzbuzz-variations/`](./exercises/02-fizzbuzz-variations/)
- **Task**: Implement an extensible rule-based FizzBuzz engine using clean switch/conditional control flow and custom delimiters.
- **Run tests**:
  ```powershell
  go test -v ./01-fundamentals/exercises/02-fizzbuzz-variations
  ```

### 3. Prime Number & Factorization Engine (Hard)
- **Path**: [`exercises/03-prime-calculator/`](./exercises/03-prime-calculator/)
- **Task**: Efficiently test primality, find primes up to N using the Sieve of Eratosthenes, and compute prime factorizations.
- **Run tests**:
  ```powershell
  go test -v ./01-fundamentals/exercises/03-prime-calculator
  ```

---

## 🧩 Challenge

Build a command-line calculator utility in [`challenges/cli-calculator/`](./challenges/cli-calculator/):
- Parse operation (`add`, `sub`, `mul`, `div`, `mod`, `pow`) and numeric operands from CLI arguments.
- Safely handle division by zero and invalid inputs with helpful error messages.
- Verify with unit tests.

---

## ⚠️ Common Mistakes

1. **Trying to use `:=` outside a function**: Short declaration `:=` is only valid inside function bodies. At package level, always use `var`.
2. **Variable Shadowing in If blocks**:
   ```go
   var x = 10
   if true {
       x := 20 // Shadows outer x! Outer x is still 10 outside this if block.
   }
   ```
3. **Mismatched Types in Arithmetic**:
   ```go
   var a int = 5
   var b float64 = 2.5
   // c := a * b  // COMPILE ERROR: invalid operation: mismatched types int and float64
   c := float64(a) * b // Correct
   ```

---

## 🔄 Review

Answer these questions to test your conceptual recall:
1. What is the zero value of a `string`, an `int`, a `bool`, and an uninitialized pointer?
2. Why does Go forbid unused variables inside functions, but allow unused package-level variables?
3. How does Go's `switch` statement differ from C/C++/Java regarding `break` and `fallthrough`?
4. What is an untyped constant and why does `const X = 42` allow assignment to both `int32` and `float64`?

---

## ✅ Completion Criteria

- [ ] All 3 exercises in `exercises/` pass unit tests (`go test -v ./...`).
- [ ] Challenge in `challenges/cli-calculator/` is functional and tested.
- [ ] Code passes `go vet` and `go fmt`.
- [ ] You can answer all 4 Review questions without referring to notes.
- [ ] Run `review` to obtain feedback from **Go Reviewer**.

---

## 🔗 Connections

- **Prerequisites**: [Stage 00: Orientation](../00-orientation/README.md) (Toolchain & Environment).
- **Next Stage**: [Stage 02: Data Structures](../02-data-structures/README.md) (Arrays, Slices, Slice Headers, Maps).
