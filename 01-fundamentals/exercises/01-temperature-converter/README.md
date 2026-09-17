# Exercise 01: Temperature Converter (Tier 1 - Easy)

## 🎯 Problem Statement
Implement conversion functions between Celsius, Fahrenheit, and Kelvin temperature scales.

The physics formulas:
* **Celsius to Fahrenheit**: `F = (C * 9/5) + 32`
* **Fahrenheit to Celsius**: `C = (F - 32) * 5/9`
* **Celsius to Kelvin**: `K = C + 273.15`
* **Kelvin to Celsius**: `C = K - 273.15`

Additionally, temperatures cannot drop below Absolute Zero (`0 Kelvin` or `-273.15 Celsius`).
If an input is below Absolute Zero, the function must return an error using `errors.New` or `fmt.Errorf`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement the functions marked with `// TODO: Implement`.
3. Do not change function signatures.
4. Run the test suite:
   ```powershell
   go test -v ./01-fundamentals/exercises/01-temperature-converter
   ```

---

## 💡 Hints

<details>
<summary>Hint 1: Integer Division Trap (Click to expand)</summary>

In Go, `9 / 5` evaluates to `1` because both operands are integers!
Make sure to use floating-point literals: `9.0 / 5.0`.
</details>

<details>
<summary>Hint 2: Error Handling (Click to expand)</summary>

Use `import "errors"` and return `errors.New("temperature below absolute zero")` when validation fails.
When returning successfully, return `nil` for the error: `return result, nil`.
</details>

<details>
<summary>Hint 3: Rounding in Tests (Click to expand)</summary>

Floating-point math can introduce minute precision differences (e.g. `99.99999999999999`). The test suite checks values within a tiny delta (`1e-4`).
</details>
