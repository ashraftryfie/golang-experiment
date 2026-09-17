# Stage 00: Orientation & Environment

## 🎯 Goal
Understand the Go runtime architecture, compilation model, CLI toolchain, project layout conventions, and how tests work before writing application code.

---

## 🧠 Core Concepts
- **Ahead-Of-Time (AOT) Compilation**: Go compiles directly to a single, standalone machine binary containing both your code and the Go runtime (including the garbage collector and scheduler).
- **The Go Toolchain**: `go run`, `go build`, `go test`, `go fmt`, and `go vet` form an all-in-one developer environment with no external build system required.
- **Strict Formatting & Opinionated Style**: Go code has one official format enforced by `gofmt`. In Go, tabs are used for indentation, not spaces. Unused variables and unused imports are compile-time errors.
- **Single Module Architecture**: `go.mod` declares the module identity and tracks direct and indirect dependencies with cryptographic checksums in `go.sum`.

---

## 🔍 Mental Model

Think of Go as a self-contained astronaut suit:

```text
┌─────────────────────────────────────────────────┐
│              Standalone Go Binary               │
│                                                 │
│   ┌─────────────────────────────────────────┐   │
│   │           Your Application Code         │   │
│   └─────────────────────────────────────────┘   │
│                                                 │
│   ┌─────────────────────────────────────────┐   │
│   │               Go Runtime                │   │
│   │  • Memory Allocator & Garbage Collector │   │
│   │  • Goroutine Scheduler (M:N)            │   │
│   │  • Channel Synchronization Engine       │   │
│   └─────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────┘
                         │
                         ▼
             OS Kernel & Hardware CPU
        (No JVM, No Node runtime, No Python VM)
```

In Python or Node.js, your code needs an interpreter installed on the target server.
In Go, your compiled binary is completely self-sufficient. You copy a single binary onto a clean Linux server or minimal container, and it executes directly on the CPU.

---

## 💻 Examples

The canonical starting point in Go:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Gopher!")
}
```

Key observations:
1. `package main` tells the Go compiler that this file produces an executable program rather than a shared library.
2. `func main()` is the entry point of the binary. It takes no arguments and returns no value.
3. Unused imports are illegal: if you import `"os"` but never use it, Go refuses to compile.

Run and verify:
```powershell
go run ./00-orientation/examples/hello
```

---

## 🧪 Exercises

### 1. Environment Verification (Easy)
- Verify your local toolchain by running:
  ```powershell
  go version
  go env GOPATH GOROOT
  ```
- Navigate to `00-orientation/examples/hello/` and execute the unit test:
  ```powershell
  go test -v ./00-orientation/examples/hello
  ```

---

## 🧩 Challenge

Run `go build` to generate a native executable binary, verify its size, and run it directly from PowerShell:
```powershell
go build -o ./00-orientation/hello.exe ./00-orientation/examples/hello
./00-orientation/hello.exe
```
Notice how fast execution is compared to interpreted scripts. Clean up the `.exe` afterwards (it is ignored by `.gitignore`).

---

## ⚠️ Common Mistakes

1. **Attempting to leave unused imports or variables**: In other languages this is a warning; in Go, it is a hard compile error. If you need a temporary placeholder, assign to the blank identifier `_ = myVar`.
2. **Opening braces on a new line**:
   ```go
   // COMPILE ERROR
   func main() 
   {
   }
   ```
   Go's lexer inserts semicolons automatically at end of lines. An opening brace MUST be on the same line as the function signature or statement.
3. **Using spaces instead of tabs**: Go strictly uses tabs for indentation. Let `go fmt` handle formatting automatically.

---

## 🔄 Review

Before moving to Stage 01, answer these questions:
1. What does `package main` signify to the Go compiler compared to `package mylib`?
2. Why doesn't Go require a virtual machine (like Java's JVM or Node's V8) on the host machine to run production binaries?
3. What is the difference between `go run` and `go build`?
4. What command automatically formats all Go files in your project?

---

## ✅ Completion Criteria

- [x] Go 1.22+ is installed and verified via `go version`.
- [x] `00-orientation/examples/hello/main.go` runs successfully.
- [x] Unit test in `00-orientation/examples/hello/main_test.go` passes cleanly (`PASS`).
- [x] You can explain the difference between `go run` and `go build`.

---

## 🔗 Connections

- **Next Stage**: [Stage 01: Go Fundamentals](../01-fundamentals/README.md) explores Go's static type system, zero values, control flow, functions, and hands-on exercises.
