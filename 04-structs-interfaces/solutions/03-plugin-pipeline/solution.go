package plugin

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
	for _, pl := range p.plugins {
		if initializer, ok := pl.(Initializer); ok {
			if err := initializer.Init(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *PluginPipeline) HealthCheckAll() map[string]bool {
	status := make(map[string]bool, len(p.plugins))
	for _, pl := range p.plugins {
		if checker, ok := pl.(HealthChecker); ok {
			status[pl.Name()] = checker.Healthy()
		} else {
			status[pl.Name()] = true // Default healthy if no checker defined
		}
	}
	return status
}

func (p *PluginPipeline) RunPipeline(input string) (string, error) {
	current := input
	for _, pl := range p.plugins {
		next, err := pl.Execute(current)
		if err != nil {
			return "", err
		}
		current = next
	}
	return current, nil
}
