package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WriteAtomic writes data to a temporary file, calls Sync, closes it, and renames it to dest.
// This guarantees crash-resilience: readers never see a half-written file.
func WriteAtomic(destPath string, data []byte) error {
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir failed: %w", err)
	}

	// Create temp file in the same directory to guarantee same filesystem for atomic rename
	tmpFile, err := os.CreateTemp(dir, "tmp-*.data")
	if err != nil {
		return fmt.Errorf("create temp failed: %w", err)
	}
	tmpName := tmpFile.Name()

	// Ensure cleanup if something fails before the rename
	defer func() {
		if tmpName != "" {
			_ = os.Remove(tmpName)
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("write failed: %w", err)
	}

	// Flush dirty pages to stable storage before rename
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("sync failed: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close failed: %w", err)
	}

	// Atomic replacement
	if err := os.Rename(tmpName, destPath); err != nil {
		return fmt.Errorf("atomic rename failed: %w", err)
	}

	tmpName = "" // Prevent defer cleanup from removing the successfully moved file
	return nil
}

func main() {
	target := filepath.Join(os.TempDir(), "golang_experiment_atomic_demo.txt")
	defer os.Remove(target)

	payload := []byte("Safe, atomic persistence in Go using os.CreateTemp and os.Rename!\n")
	if err := WriteAtomic(target, payload); err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}

	// Inspect existence & metadata using os.Stat
	info, err := os.Stat(target)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println("File does not exist!")
		return
	} else if err != nil {
		fmt.Printf("Stat error: %v\n", err)
		return
	}

	fmt.Printf("Saved atomically: %s (Size: %d bytes, Mode: %v)\n", info.Name(), info.Size(), info.Mode())

	// Read content
	f, err := os.Open(target)
	if err != nil {
		fmt.Printf("Open error: %v\n", err)
		return
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("Read error: %v\n", err)
		return
	}

	fmt.Printf("File content: %s", string(content))
}
