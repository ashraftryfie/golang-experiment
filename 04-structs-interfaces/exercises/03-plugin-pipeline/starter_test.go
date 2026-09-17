package plugin

import (
	"errors"
	"strings"
	"testing"
)

type UpperPlugin struct{}

func (u UpperPlugin) Name() string { return "upper" }
func (u UpperPlugin) Execute(p string) (string, error) {
	return strings.ToUpper(p), nil
}

type ExclamationPlugin struct {
	initialized bool
}

func (e *ExclamationPlugin) Name() string { return "exclamation" }
func (e *ExclamationPlugin) Execute(p string) (string, error) {
	return p + "!", nil
}
func (e *ExclamationPlugin) Init() error {
	e.initialized = true
	return nil
}
func (e *ExclamationPlugin) Healthy() bool {
	return e.initialized
}

func TestPluginPipeline(t *testing.T) {
	pipe := NewPipeline()
	pipe.Register(UpperPlugin{})
	excl := &ExclamationPlugin{}
	pipe.Register(excl)

	err := pipe.InitAll()
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: InitAll is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if !excl.initialized {
		t.Errorf("expected exclamation plugin to be initialized")
	}

	health := pipe.HealthCheckAll()
	if !health["upper"] || !health["exclamation"] {
		t.Errorf("expected both plugins to be healthy, got %v", health)
	}

	out, err := pipe.RunPipeline("hello")
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if out != "HELLO!" {
		t.Errorf("RunPipeline('hello') = %q, want 'HELLO!'", out)
	}
}
