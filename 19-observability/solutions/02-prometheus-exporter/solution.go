package promexporter

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

var (
	ErrDuplicateMetric = errors.New("metric with this name already registered")
)

type metricEntry struct {
	name      string
	help      string
	mType     string
	labelKeys []string
	values    map[string]float64
	mu        sync.RWMutex
}

func (m *metricEntry) labelKey(vals []string) string {
	if len(m.labelKeys) == 0 || len(vals) == 0 {
		return ""
	}
	pairs := make([]string, len(m.labelKeys))
	for i, k := range m.labelKeys {
		v := ""
		if i < len(vals) {
			v = vals[i]
		}
		pairs[i] = fmt.Sprintf(`%s="%s"`, k, v)
	}
	return "{" + strings.Join(pairs, ",") + "}"
}

type Counter struct {
	entry *metricEntry
}

func (c *Counter) Inc(labelValues ...string) {
	c.Add(1, labelValues...)
}

func (c *Counter) Add(val float64, labelValues ...string) {
	if val < 0 {
		return
	}
	key := c.entry.labelKey(labelValues)
	c.entry.mu.Lock()
	defer c.entry.mu.Unlock()
	c.entry.values[key] += val
}

type Gauge struct {
	entry *metricEntry
}

func (g *Gauge) Set(val float64, labelValues ...string) {
	key := g.entry.labelKey(labelValues)
	g.entry.mu.Lock()
	defer g.entry.mu.Unlock()
	g.entry.values[key] = val
}

func (g *Gauge) Inc(labelValues ...string) {
	g.Add(1, labelValues...)
}

func (g *Gauge) Dec(labelValues ...string) {
	g.Add(-1, labelValues...)
}

func (g *Gauge) Add(delta float64, labelValues ...string) {
	key := g.entry.labelKey(labelValues)
	g.entry.mu.Lock()
	defer g.entry.mu.Unlock()
	g.entry.values[key] += delta
}

type Registry struct {
	metrics map[string]*metricEntry
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		metrics: make(map[string]*metricEntry),
	}
}

func (r *Registry) RegisterCounter(name, help string, labelKeys []string) (*Counter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.metrics[name]; exists {
		return nil, ErrDuplicateMetric
	}

	entry := &metricEntry{
		name:      name,
		help:      help,
		mType:     "counter",
		labelKeys: labelKeys,
		values:    make(map[string]float64),
	}
	r.metrics[name] = entry
	return &Counter{entry: entry}, nil
}

func (r *Registry) RegisterGauge(name, help string, labelKeys []string) (*Gauge, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.metrics[name]; exists {
		return nil, ErrDuplicateMetric
	}

	entry := &metricEntry{
		name:      name,
		help:      help,
		mType:     "gauge",
		labelKeys: labelKeys,
		values:    make(map[string]float64),
	}
	r.metrics[name] = entry
	return &Gauge{entry: entry}, nil
}

func (r *Registry) Export() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var buf bytes.Buffer

	// Sort metric names for deterministic output
	names := make([]string, 0, len(r.metrics))
	for n := range r.metrics {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		m := r.metrics[name]
		m.mu.RLock()

		buf.WriteString(fmt.Sprintf("# HELP %s %s\n", m.name, m.help))
		buf.WriteString(fmt.Sprintf("# TYPE %s %s\n", m.name, m.mType))

		if len(m.values) == 0 {
			buf.WriteString(fmt.Sprintf("%s 0\n", m.name))
		} else {
			// Sort label combinations
			keys := make([]string, 0, len(m.values))
			for k := range m.values {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			for _, k := range keys {
				val := m.values[k]
				// Format clean float/int representation
				if val == float64(int64(val)) {
					buf.WriteString(fmt.Sprintf("%s%s %d\n", m.name, k, int64(val)))
				} else {
					buf.WriteString(fmt.Sprintf("%s%s %g\n", m.name, k, val))
				}
			}
		}
		buf.WriteString("\n")
		m.mu.RUnlock()
	}

	return buf.String()
}
