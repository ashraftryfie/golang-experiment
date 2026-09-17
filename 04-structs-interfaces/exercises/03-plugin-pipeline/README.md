# Exercise 03: Plugin Lifecycle Pipeline (Tier 3 - Hard)

## 🎯 Problem Statement
Implement a flexible plugin lifecycle pipeline using interface segregation and runtime capability detection via type assertions.

### Core Interface
Every plugin must implement `Plugin`:
```go
type Plugin interface {
    Name() string
    Execute(payload string) (string, error)
}
```

### Optional Capability Interfaces
A plugin may optionally implement one or both of these interfaces:
```go
// Initializer is checked before pipeline execution.
type Initializer interface {
    Init() error
}

// HealthChecker is checked to assess plugin readiness.
type HealthChecker interface {
    Healthy() bool
}
```

### Pipeline Requirements
`PluginPipeline`:
- `Register(plugin Plugin)`: Registers a plugin into the pipeline.
- `InitAll() error`: Iterates over registered plugins; if a plugin satisfies `Initializer`, calls its `Init()` method. Stops and returns error if initialization fails.
- `HealthCheckAll() map[string]bool`: Returns map of plugin names to health status. If a plugin does NOT implement `HealthChecker`, it is assumed healthy (`true`).
- `RunPipeline(input string) (string, error)`: Passes `input` sequentially through all plugins in registration order, returning the final transformed output.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Use type assertions (`if init, ok := p.(Initializer); ok { ... }`).
3. Run tests:
   ```powershell
   go test -v ./04-structs-interfaces/exercises/03-plugin-pipeline
   ```
