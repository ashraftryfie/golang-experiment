package main

import (
	"fmt"
	"sort"
)

type LogRecord struct {
	Service    string
	StatusCode int
	LatencyMs  float64
}

type ServiceMetrics struct {
	TotalRequests int
	ErrorCount    int
	AvgLatencyMs  float64
	MaxLatencyMs  float64
}

type serviceAccumulator struct {
	totalRequests int
	errorCount    int
	totalLatency  float64
	maxLatency    float64
}

type ServiceSummary struct {
	Service    string
	ErrorCount int
}

type Aggregator struct {
	data map[string]*serviceAccumulator
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		data: make(map[string]*serviceAccumulator),
	}
}

func (a *Aggregator) Ingest(record LogRecord) {
	acc, exists := a.data[record.Service]
	if !exists {
		acc = &serviceAccumulator{}
		a.data[record.Service] = acc
	}

	acc.totalRequests++
	if record.StatusCode >= 500 {
		acc.errorCount++
	}
	acc.totalLatency += record.LatencyMs
	if record.LatencyMs > acc.maxLatency {
		acc.maxLatency = record.LatencyMs
	}
}

func (a *Aggregator) GetMetrics(service string) (ServiceMetrics, bool) {
	acc, exists := a.data[service]
	if !exists || acc.totalRequests == 0 {
		return ServiceMetrics{}, false
	}

	return ServiceMetrics{
		TotalRequests: acc.totalRequests,
		ErrorCount:    acc.errorCount,
		AvgLatencyMs:  acc.totalLatency / float64(acc.totalRequests),
		MaxLatencyMs:  acc.maxLatency,
	}, true
}

func (a *Aggregator) Services() []string {
	services := make([]string, 0, len(a.data))
	for s := range a.data {
		services = append(services, s)
	}
	sort.Strings(services)
	return services
}

func (a *Aggregator) TopFailingServices(limit int) []ServiceSummary {
	if limit <= 0 {
		return []ServiceSummary{}
	}

	summaries := make([]ServiceSummary, 0, len(a.data))
	for name, acc := range a.data {
		summaries = append(summaries, ServiceSummary{
			Service:    name,
			ErrorCount: acc.errorCount,
		})
	}

	sort.Slice(summaries, func(i, j int) bool {
		if summaries[i].ErrorCount != summaries[j].ErrorCount {
			return summaries[i].ErrorCount > summaries[j].ErrorCount
		}
		return summaries[i].Service < summaries[j].Service
	})

	if limit > len(summaries) {
		limit = len(summaries)
	}
	return summaries[:limit]
}

func main() {
	agg := NewAggregator()
	agg.Ingest(LogRecord{Service: "auth", StatusCode: 200, LatencyMs: 45.2})
	agg.Ingest(LogRecord{Service: "auth", StatusCode: 500, LatencyMs: 120.0})
	agg.Ingest(LogRecord{Service: "billing", StatusCode: 200, LatencyMs: 88.5})

	for _, s := range agg.Services() {
		m, _ := agg.GetMetrics(s)
		fmt.Printf("[%s] Total: %d, Errors: %d, Avg: %.2fms, Max: %.2fms\n",
			s, m.TotalRequests, m.ErrorCount, m.AvgLatencyMs, m.MaxLatencyMs)
	}
}
