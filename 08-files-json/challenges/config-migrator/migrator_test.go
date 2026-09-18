package configmigrator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadOrMigrateV1ToV2(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "app.json")

	// Write initial v1 config
	v1Raw := `{
		"version": 1,
		"host": "0.0.0.0",
		"port": 9090
	}`
	if err := os.WriteFile(configPath, []byte(v1Raw), 0644); err != nil {
		t.Fatalf("failed to write initial v1 config: %v", err)
	}

	// Run migration
	cfg, err := LoadOrMigrate(configPath)
	if err != nil {
		t.Fatalf("LoadOrMigrate failed: %v", err)
	}

	if cfg.Version != 2 {
		t.Errorf("expected version 2, got %d", cfg.Version)
	}
	if cfg.Server.Host != "0.0.0.0" || cfg.Server.Port != 9090 {
		t.Errorf("server config mismatch: %+v", cfg.Server)
	}
	if cfg.TLS != false {
		t.Errorf("expected TLS to default to false, got %v", cfg.TLS)
	}

	// Verify the file on disk was replaced with v2 JSON
	diskData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read migrated file: %v", err)
	}

	var diskCfg ConfigV2
	if err := json.Unmarshal(diskData, &diskCfg); err != nil {
		t.Fatalf("failed to parse disk JSON: %v", err)
	}
	if diskCfg.Version != 2 || diskCfg.Server.Port != 9090 {
		t.Errorf("disk config mismatch: %+v", diskCfg)
	}

	// Loading again should load v2 directly without error
	cfg2, err := LoadOrMigrate(configPath)
	if err != nil {
		t.Fatalf("subsequent LoadOrMigrate failed: %v", err)
	}
	if cfg2.Version != 2 {
		t.Errorf("subsequent load version = %d, want 2", cfg2.Version)
	}
}

func TestLoadOrMigrateNotExist(t *testing.T) {
	nonExistent := filepath.Join(t.TempDir(), "does_not_exist.json")
	_, err := LoadOrMigrate(nonExistent)
	if err == nil {
		t.Fatalf("expected error for non-existent file, got nil")
	}
}
