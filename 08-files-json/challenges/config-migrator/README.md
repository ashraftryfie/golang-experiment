# Challenge: Atomic Configuration Loader & Schema Migrator

## 🎯 Objective
Build a production-grade configuration manager that loads JSON configuration from disk, applies schema migrations (e.g. migrating v1 flat layout to v2 nested structure), and writes back updates using crash-safe atomic file replacement.

## 📋 Requirements
1. **Schema Definitions**:
   - `ConfigV1`:
     ```go
     type ConfigV1 struct {
         Version int    `json:"version"`
         Host    string `json:"host"`
         Port    int    `json:"port"`
     }
     ```
   - `ConfigV2`:
     ```go
     type ServerConfig struct {
         Host string `json:"host"`
         Port int    `json:"port"`
     }
     type ConfigV2 struct {
         Version int          `json:"version"`
         Server  ServerConfig `json:"server"`
         TLS     bool         `json:"tls"`
     }
     ```
2. **Operations**:
   - `LoadOrMigrate(configPath string) (*ConfigV2, error)`:
     - Inspects file existence using `os.Stat(configPath)` and `errors.Is(err, os.ErrNotExist)`.
     - Detects schema version by inspecting `"version"`.
     - If v1, transforms to v2 with `TLS: false`, writes v2 atomically back to `configPath`, and returns the v2 struct.
     - If v2, decodes and returns directly.
3. **Crash-Safe Atomic Write**:
   - Uses `os.CreateTemp` in the same directory, writes JSON with indent, flushes dirty pages with `f.Sync()`, closes, and renames atomically with `os.Rename`.

## 🧪 Verification
```powershell
go test -v ./08-files-json/challenges/config-migrator/...
```
