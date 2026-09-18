package configmigrator

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ConfigV1 struct {
	Version int    `json:"version"`
	Host    string `json:"host"`
	Port    int    `json:"port"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ConfigV2 struct {
	Version int          `json:"version"`
	Server  ServerConfig `json:"server"`
	TLS     bool         `json:"tls"`
}

// VersionHeader is used to inspect the schema version before full deserialization.
type VersionHeader struct {
	Version int `json:"version"`
}

// SaveAtomic writes cfg to filePath atomically using a temporary file and sync.
func SaveAtomic(filePath string, cfg any) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	tmpFile, err := os.CreateTemp(dir, "cfg-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file failed: %w", err)
	}
	tmpName := tmpFile.Name()

	defer func() {
		if tmpName != "" {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("write temp failed: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("sync temp failed: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp failed: %w", err)
	}

	if err := os.Rename(tmpName, filePath); err != nil {
		return fmt.Errorf("atomic rename failed: %w", err)
	}

	tmpName = "" // Prevent defer from deleting moved file
	return nil
}

// LoadOrMigrate reads the configuration file from disk.
// If it is v1, it automatically upgrades it to v2 and persists it atomically.
func LoadOrMigrate(filePath string) (*ConfigV2, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("config file does not exist: %w", err)
		}
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	var header VersionHeader
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("parse version header failed: %w", err)
	}

	switch header.Version {
	case 1:
		var v1 ConfigV1
		if err := json.Unmarshal(data, &v1); err != nil {
			return nil, fmt.Errorf("parse v1 config failed: %w", err)
		}

		// Migrate to v2
		v2 := &ConfigV2{
			Version: 2,
			Server: ServerConfig{
				Host: v1.Host,
				Port: v1.Port,
			},
			TLS: false,
		}

		// Persist migrated config atomically
		if err := SaveAtomic(filePath, v2); err != nil {
			return nil, fmt.Errorf("save migrated config failed: %w", err)
		}

		return v2, nil

	case 2:
		var v2 ConfigV2
		if err := json.Unmarshal(data, &v2); err != nil {
			return nil, fmt.Errorf("parse v2 config failed: %w", err)
		}
		return &v2, nil

	default:
		return nil, fmt.Errorf("unsupported config schema version: %d", header.Version)
	}
}
