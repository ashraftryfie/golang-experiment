# Exercise 01: Functional Transformation Pipelines (Tier 1 - Easy)

## 🎯 Problem Statement
Implement the classic functional programming trio for integer slices using first-class functions:

1. `Map(nums []int, fn func(int) int) []int`:
   - Returns a new slice with `fn` applied to each element.
2. `Filter(nums []int, predicate func(int) bool) []int`:
   - Returns a new slice with only elements satisfying `predicate(n) == true`.
3. `Reduce(nums []int, initial int, accumulator func(acc, current int) int) int`:
   - Combines all elements into a single aggregate value starting from `initial`.
4. `Compose(f, g func(int) int) func(int) int`:
   - Returns a new function representing mathematical composition: `f(g(x))`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement the functions according to docstrings.
3. Run tests:
   ```powershell
   go test -v ./03-functions-methods/exercises/01-function-pipelines
   ```
