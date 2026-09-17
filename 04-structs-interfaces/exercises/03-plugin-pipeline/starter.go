package plugin

import (
	"errors"
)

var ErrNotImplemented = errors.New("TODO: implement")

type Plugin interface {
	Name() string
	Execute(payload string) (string, error)
}

type Initializer interface {
	Init() error
}

type HealthChecker interface {
	Healthy() bool
}

type PluginPipeline struct {
	plugins []Plugin
}

func NewPipeline() *PluginPipeline {
	return &PluginPipeline{
		plugins: make([]Plugin, 0),
	}
}

func (p *PluginPipeline) Register(pl Plugin) {
	p.plugins = append(p.plugins, pl)
}

func (p *PluginPipeline) InitAll() error {
	// TODO: For each plugin, check if it implements Initializer via type assertion and call Init()
	return ErrNotImplemented
}

func (p *PluginPipeline) HealthCheckAll() map[string]bool {
	// TODO: For each plugin, check if it implements HealthChecker via type assertion
	return nil
}

func (p *PluginPipeline) RunPipeline(input string) (string, error) {
	// TODO: Pipe input through all plugins sequentially
	return "", ErrNotImplemented
}
