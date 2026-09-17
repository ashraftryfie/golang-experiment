# Go Toolchain Commands Cheat Sheet

Essential commands for building, testing, linting, and maintaining Go modules.

---

## 1. Running & Building

| Command | Action | Notes |
| :--- | :--- | :--- |
| `go run .` | Compile and run immediate package in memory | Fast iteration for scripts/binaries |
| `go build -o bin/app ./cmd/app` | Compile binary into destination directory | Produces native OS executable |
| `CGO_ENABLED=0 go build -ldflags="-s -w"` | Compile fully static binary without C dependencies | Strips debug symbols (`-s`) and DWARF (`-w`) for smallest binary |
| `go clean -cache` | Clear Go compiler build cache | Useful if encountering corrupted cache |

---

## 2. Testing & Profiling

| Command | Action | Notes |
| :--- | :--- | :--- |
| `go test ./...` | Run all tests across the entire module | Recursively executes all `*_test.go` |
| `go test -v ./...` | Run tests with verbose output | Prints names and results of every subtest |
| `go test -race ./...` | Run tests with Go runtime race detector | **Critical for concurrent code** |
| `go test -run TestName ./path` | Run a specific test matching regex pattern | Saves time during focused debugging |
| `go test -bench=. -benchmem` | Run benchmarks and report memory allocations | Measures ns/op and B/op |
| `go test -fuzz=FuzzName` | Run fuzz tests with randomized inputs | Finds edge case panics and memory bugs |
| `go test -coverprofile=c.out ./...` | Generate code coverage profile | View in browser with `go tool cover -html=c.out` |

---

## 3. Formatting & Static Analysis

| Command | Action | Notes |
| :--- | :--- | :--- |
| `go fmt ./...` | Format all Go code according to standard style | Standard indentation (tabs) and spacing |
| `gofmt -s -w .` | Format and simplify code in-place | Replaces verbose slice slices and loops |
| `go vet ./...` | Run Go compiler static analysis | Catches unreachable code, printf format mismatches, shadow variables |

---

## 4. Module Management

| Command | Action | Notes |
| :--- | :--- | :--- |
| `go mod init <module-path>` | Initialize a new `go.mod` file | Sets root package namespace |
| `go mod tidy` | Add missing module dependencies and remove unused ones | Cleans up `go.mod` and `go.sum` |
| `go mod download` | Download dependencies to local cache | Caches in `$GOPATH/pkg/mod` |
| `go mod verify` | Verify cryptographic hashes of dependencies in `go.sum` | Ensures supply chain integrity |
