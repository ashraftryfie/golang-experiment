# Challenge: CLI Calculator Utility

## 🎯 Goal
Build a functional, tested command-line calculator utility that takes operands and an operation from command-line arguments, calculates the result, and prints formatted output or descriptive errors.

---

## 📋 Requirements

### CLI Invocation Syntax
```powershell
go run ./01-fundamentals/challenges/cli-calculator <operation> <operand1> <operand2>
```

### Supported Operations
- `add`: `operand1 + operand2`
- `sub`: `operand1 - operand2`
- `mul`: `operand1 * operand2`
- `div`: `operand1 / operand2` (must handle division by zero error)
- `mod`: integer modulo `int(operand1) % int(operand2)`
- `pow`: `math.Pow(operand1, operand2)`

### Error Handling Requirements
- Must print usage help if fewer than 3 arguments are passed.
- Must return an error if non-numeric values are passed as operands.
- Must return `"error: division by zero is undefined"` on divide by zero.
- Must return `"error: unsupported operation <op>"` for unknown operations.

---

## 🧪 Testing
```powershell
go test -v ./01-fundamentals/challenges/cli-calculator
```
