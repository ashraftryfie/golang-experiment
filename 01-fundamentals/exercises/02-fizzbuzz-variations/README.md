# Exercise 02: FizzBuzz Variations (Tier 2 - Medium)

## 🎯 Problem Statement
Implement a flexible, rule-based FizzBuzz engine that evaluates integer sequences according to configurable divisor rules, returning formatted slices of strings.

### Requirements
1. `ClassicFizzBuzz(n int) ([]string, error)`
   - Given an integer `n >= 1`, return a slice of strings from `1` to `n`.
   - Numbers divisible by 3 become `"Fizz"`.
   - Numbers divisible by 5 become `"Buzz"`.
   - Numbers divisible by both 3 and 5 become `"FizzBuzz"`.
   - All other numbers are represented as their base-10 string representation (e.g. `"1"`, `"2"`).
   - If `n < 1`, return an error: `errors.New("n must be at least 1")`.

2. `CustomFizzBuzz(n int, rules []Rule) ([]string, error)`
   - Takes a slice of `Rule` structs:
     ```go
     type Rule struct {
         Divisor int
         Word    string
     }
     ```
   - For each number `1..n`, check rules in order. If a number is divisible by a rule's divisor, append the rule's word.
   - If multiple rules match (e.g. Divisor 3 "Fizz", Divisor 7 "Bazz"), concatenate the words in rule order (e.g. `"FizzBazz"`).
   - If no rules match, return the number string.
   - If `n < 1` or if any `Rule.Divisor <= 0`, return an appropriate error.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `ClassicFizzBuzz` and `CustomFizzBuzz`.
3. Use `strconv.Itoa()` to convert integers to strings.
4. Run tests:
   ```powershell
   go test -v ./01-fundamentals/exercises/02-fizzbuzz-variations
   ```

---

## 💡 Hints

<details>
<summary>Hint 1: Integer to String (Click to expand)</summary>

Import `"strconv"` and use `strconv.Itoa(i)` rather than `fmt.Sprintf("%d", i)` for better efficiency and simplicity.
</details>

<details>
<summary>Hint 2: Pre-allocating Slices (Click to expand)</summary>

Since you know the result will contain exactly `n` elements, allocate with `make([]string, 0, n)` to avoid slice re-allocations as you append!
</details>
