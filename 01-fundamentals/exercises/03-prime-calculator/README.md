# Exercise 03: Prime Calculator & Factorization (Tier 3 - Hard)

## 🎯 Problem Statement
Implement an efficient prime calculation engine that handles:
1. **Primality Testing (`IsPrime(n int) bool`)**: Determines if an integer `n` is prime in `O(sqrt(N))` time.
2. **Sieve of Eratosthenes (`SieveOfEratosthenes(limit int) []int`)**: Efficiently generates all prime numbers up to `limit`.
3. **Prime Factorization (`PrimeFactors(n int) ([]int, error)`)**: Decomposes an integer `n >= 2` into its ascending prime factors (e.g. `12 -> [2, 2, 3]`).

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement the functions according to docstrings.
3. Handle invalid inputs (`n < 2`) with descriptive errors.
4. Run tests:
   ```powershell
   go test -v ./01-fundamentals/exercises/03-prime-calculator
   ```

---

## 💡 Hints

<details>
<summary>Hint 1: Primality Edge Cases (Click to expand)</summary>

Numbers `<= 1` are not prime. `2` is the only even prime. Any even number `> 2` can be immediately discarded.
You only need to check odd divisors up to `int(math.Sqrt(float64(n)))`.
</details>

<details>
<summary>Hint 2: Boolean Sieve (Click to expand)</summary>

Use a boolean slice `isComposite := make([]bool, limit+1)`.
For each prime `p`, mark multiples `p*p, p*p + p, ...` as true.
</details>
