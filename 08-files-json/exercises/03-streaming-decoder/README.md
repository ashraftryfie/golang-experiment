# Exercise 03: Streaming JSON Array Decoder & Encoder (Tier 3 - Hard)

## 🎯 Problem Statement
In microservices handling multi-gigabyte log dumps or audit streams, unmarshaling the entire JSON array with `json.Unmarshal` leads to out-of-memory (OOM) fatal crashes.

Implement an incremental streaming filter:
```go
type Transaction struct {
    ID     string  `json:"id"`
    Amount float64 `json:"amount"`
    Sender string  `json:"sender"`
}

func FilterHighValue(dst io.Writer, src io.Reader, minAmount float64) (int, error)
```

### Requirements
1. Consume the opening array delimiter `[` using `decoder.Token()`.
2. Iterate through items using `decoder.More()` and `decoder.Decode(&tx)`.
3. Filter transactions where `tx.Amount >= minAmount`.
4. Stream matching transactions to `dst` as a valid JSON array `[...]`.
5. Consume the closing bracket `]` using `decoder.Token()`.
6. Return the total count of matched transactions written.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `FilterHighValue`.
3. Run tests:
   ```powershell
   go test -v ./08-files-json/exercises/03-streaming-decoder
   ```
