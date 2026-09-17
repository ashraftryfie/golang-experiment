# Exercise 01: Slice Operations (Tier 1 - Easy)

## 🎯 Problem Statement
Implement high-performance slice manipulation functions using Go's slice slicing syntax and in-place algorithms:

1. `Deduplicate(nums []int) []int`:
   - Returns a slice containing only unique integers, preserving original order.
   - Must handle empty and nil slices gracefully.
2. `FilterEven(nums []int) []int`:
   - Filters out all odd numbers, returning only even numbers.
   - **Performance requirement**: Must reuse the backing array of the input slice in-place (allocation-free: `nums[:0]`).
3. `RemoveAtIndex(nums []int, index int) ([]int, error)`:
   - Removes the element at `index` while preserving the order of the remaining elements.
   - Returns `errors.New("index out of range")` if `index < 0` or `index >= len(nums)`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement functions marked with `TODO`.
3. Run tests:
   ```powershell
   go test -v ./02-data-structures/exercises/01-slice-ops
   ```

---

## 💡 Hints

<details>
<summary>Hint 1: In-Place Filtering Idiom (Click to expand)</summary>

You can filter in-place without allocating a new backing array:
```go
out := nums[:0]
for _, x := range nums {
    if isEven(x) {
        out = append(out, x)
    }
}
return out
```
</details>

<details>
<summary>Hint 2: Removing at Index (Click to expand)</summary>

Use Go's sub-slice append idiom:
```go
return append(nums[:index], nums[index+1:]...), nil
```
</details>
