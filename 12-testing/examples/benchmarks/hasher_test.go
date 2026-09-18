package benchmarks

import (
	"strings"
	"testing"
)

// ConcatPlus joins strings using traditional + operator (allocates each time)
func ConcatPlus(parts []string) string {
	var s string
	for _, p := range parts {
		s += p
	}
	return s
}

// ConcatBuilder joins strings using strings.Builder with pre-allocated capacity
func ConcatBuilder(parts []string) string {
	totalLen := 0
	for _, p := range parts {
		totalLen += len(p)
	}

	var sb strings.Builder
	sb.Grow(totalLen)
	for _, p := range parts {
		sb.WriteString(p)
	}
	return sb.String()
}

var benchmarkParts = []string{
	"microservice", "event", "pipeline", "streaming", "architecture",
	"production", "kubernetes", "distroless", "observability", "metrics",
}

func BenchmarkConcatPlus(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ConcatPlus(benchmarkParts)
	}
}

func BenchmarkConcatBuilder(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ConcatBuilder(benchmarkParts)
	}
}
