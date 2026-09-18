# Exercise 01: Table-Driven Testing Suite (Tier 1 - Easy)

## 🎯 Problem Statement
A URL slugifier transforms arbitrary human strings into clean, SEO-friendly URL segments:
```go
func Slugify(s string) string
```

### Transformation Rules
1. Converts all letters to lowercase.
2. Replaces spaces and underscores with hyphens `-`.
3. Strips out non-alphanumeric characters (except hyphens).
4. Collapses consecutive hyphens (`---`) into a single hyphen (`-`).
5. Trims leading and trailing hyphens.

Implement a comprehensive table-driven test suite covering standard and edge cases (empty strings, all symbols, consecutive hyphens, mixed casing).

## 🛠️ Instructions
1. Open [`starter_test.go`](./starter_test.go).
2. Complete the test cases and helper assertions.
3. Run tests:
   ```powershell
   go test -v ./12-testing/exercises/01-table-suite
   ```
