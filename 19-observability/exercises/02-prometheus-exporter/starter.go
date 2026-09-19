package promexporter

import (
	"errors"
)

var (
	ErrNotImplemented  = errors.New("exercise not implemented yet")
	ErrDuplicateMetric = errors.New("metric with this name already registered")
)

type Counter struct{}

func (c *Counter) Inc(labelValues ...string)              {}
func (c *Counter) Add(val float64, labelValues ...string) {}

type Gauge struct{}

func (g *Gauge) Set(val float64, labelValues ...string) {}
func (g *Gauge) Inc(labelValues ...string)              {}
func (g *Gauge) Dec(labelValues ...string)              {}

type Registry struct{}

func NewRegistry() *Registry {
	return nil
}

func (r *Registry) RegisterCounter(name, help string, labelKeys []string) (*Counter, error) {
	return nil, ErrNotImplemented
}

func (r *Registry) RegisterGauge(name, help string, labelKeys []string) (*Gauge, error) {
	return nil, ErrNotImplemented
}

func (r *Registry) Export() string {
	return ""
}
