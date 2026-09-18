# Challenge: Resilient Key-Value Parser with Full Test Suite

## 🎯 Objective
Build a robust configuration string parser that extracts key-value pairs separated by `&` or `;`, ignores comments (`#`), trims whitespace, and is verified with table-driven tests, allocation benchmarks, and native fuzzing.

## 📋 Requirements
1. **Parser Signature**:
   ```go
   func ParseKV(input string) (map[string]string, error)
   ```
2. **Behavior**:
   - Supports pairs like `host=localhost&port=8080;env=prod`.
   - Skips empty pairs (e.g. `&&` or `;;`).
   - Lines or segments starting with `#` are treated as comments and ignored.
   - Trims whitespace around keys and values.
   - If a pair has no `=` delimiter, returns a descriptive error: `missing '=' delimiter in pair`.
   - Empty input returns an empty map (`map[string]string{}`).
3. **Testing Suite**:
   - **Table-Driven Tests**: Validate standard formats, edge cases, whitespace, comments, and invalid inputs.
   - **Allocation Benchmark**: `BenchmarkParseKV(b *testing.B)` with `b.ReportAllocs()`.
   - **Native Fuzz Test**: `FuzzParseKV(f *testing.F)` verifying the parser never panics on arbitrary string data.

## 🧪 Verification
```powershell
go test -v ./12-testing/challenges/resilient-parser/...
go test -bench=. -benchmem ./12-testing/challenges/resilient-parser/...
```
