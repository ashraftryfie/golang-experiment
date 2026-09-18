# Exercise 02: Custom `json.Marshaler` & `json.Unmarshaler` (Tier 2 - Medium)

## 🎯 Problem Statement
By default, Go's `time.Time` marshals to an RFC 3339 formatted string (e.g. `"2024-03-15T12:00:00Z"`).
Many legacy APIs or high-throughput analytics pipelines expect integer Unix timestamps in seconds (e.g. `1710504000`).

Implement a custom time type `EpochTime`:
```go
type EpochTime struct {
    time.Time
}

func (t EpochTime) MarshalJSON() ([]byte, error)
func (t *EpochTime) UnmarshalJSON(data []byte) error
```

### Requirements
1. `MarshalJSON()` must return the time formatted as an integer number of seconds: `fmt.Sprintf("%d", t.Unix())`.
2. `UnmarshalJSON(data []byte)` must parse the byte slice into an `int64` (e.g., using `strconv.ParseInt(string(data), 10, 64)`) and assign `t.Time = time.Unix(sec, 0)`.
3. Handle null or empty byte slices gracefully.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `MarshalJSON` and `UnmarshalJSON`.
3. Run tests:
   ```powershell
   go test -v ./08-files-json/exercises/02-custom-marshaler
   ```
