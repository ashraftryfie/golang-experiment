package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type MetricRecord struct {
	Timestamp int64   `json:"timestamp"`
	Service   string  `json:"service"`
	Value     float64 `json:"value"`
}

func main() {
	// Simulated large JSON array stream
	rawJSON := `[
		{"timestamp": 1700000001, "service": "auth", "value": 0.99},
		{"timestamp": 1700000002, "service": "payment", "value": 1.45},
		{"timestamp": 1700000003, "service": "order", "value": 0.32}
	]`

	reader := strings.NewReader(rawJSON)
	decoder := json.NewDecoder(reader)

	// Consume opening bracket '['
	t, err := decoder.Token()
	if err != nil {
		fmt.Printf("Token error: %v\n", err)
		return
	}
	fmt.Printf("Streaming started with token: %v\n", t)

	count := 0
	var sum float64

	// Stream elements one by one without allocating the whole array in memory
	for decoder.More() {
		var record MetricRecord
		if err := decoder.Decode(&record); err != nil {
			fmt.Printf("Decode error: %v\n", err)
			return
		}
		count++
		sum += record.Value
		fmt.Printf("Streamed Record %d: Service=%s, Value=%.2f\n", count, record.Service, record.Value)
	}

	// Consume closing bracket ']'
	tEnd, err := decoder.Token()
	if err != nil && !errorsIsEOF(err) {
		fmt.Printf("Closing token error: %v\n", err)
		return
	}
	fmt.Printf("Stream completed token: %v | Total: %d, Average: %.2f\n", tEnd, count, sum/float64(count))
}

func errorsIsEOF(err error) bool {
	return err == io.EOF
}
