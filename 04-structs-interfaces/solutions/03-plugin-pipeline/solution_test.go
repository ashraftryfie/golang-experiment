package plugin

import (
	"strings"
	"testing"
)

type PrefixPlugin struct {
	prefix string
}

func (p PrefixPlugin) Name() string { return "prefix" }
func (p PrefixPlugin) Execute(s string) (string, error) {
	return p.prefix + s, nil
}

type UpperPlugin struct{}

func (u UpperPlugin) Name() string { return "upper" }
func (u UpperPlugin) Execute(s string) (string, error) {
	return strings.ToUpper(s), nil
}

func TestSolutionPluginPipeline(t *testing.T) {
	pipe := NewPipeline()
	pipe.Register(PrefixPlugin{prefix: "gopher: "})
	pipe.Register(UpperPlugin{})

	out, err := pipe.RunPipeline("speed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if out != "GOPHER: SPEED" {
		t.Errorf("got %q, want 'GOPHER: SPEED'", out)
	}

	health := pipe.HealthCheckAll()
	if !health["prefix"] || !health["upper"] {
		t.Errorf("expected all plugins to be healthy by default")
	}
}
