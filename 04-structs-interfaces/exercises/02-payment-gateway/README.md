# Exercise 02: Multi-Provider Payment Gateway (Tier 2 - Medium)

## 🎯 Problem Statement
Implement a polymorphic payment processing engine using Go interfaces:

```go
type PaymentResult struct {
    TransactionID string
    Success       bool
    Message       string
}

type PaymentProcessor interface {
    ProcessPayment(amount float64) (PaymentResult, error)
}
```

### Processors to Implement
1. `CreditCardProcessor`:
   - Validates 16-digit card number.
   - Charges amount.
2. `PayPalProcessor`:
   - Validates email address (must contain `@`).
   - Charges amount.
3. `CryptoProcessor`:
   - Validates wallet address (must start with `0x`).
   - Charges amount in crypto equivalent.
4. `PaymentGateway`:
   - Manages a map of provider names to `PaymentProcessor`.
   - Dispatches `Pay(provider string, amount float64)` to the correct registered processor.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement processors and gateway.
3. Run tests:
   ```powershell
   go test -v ./04-structs-interfaces/exercises/02-payment-gateway
   ```
