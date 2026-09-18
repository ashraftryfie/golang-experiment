package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrConfigMissing = errors.New("configuration file is missing")

func loadConfigFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Wrap both the root os.ErrNotExist and our domain sentinel ErrConfigMissing
			return nil, fmt.Errorf("app config error (%w): %w", ErrConfigMissing, err)
		}
		return nil, fmt.Errorf("reading file %s: %w", path, err)
	}
	return data, nil
}

func main() {
	_, err := loadConfigFile("non_existent_config.yaml")
	if err != nil {
		fmt.Printf("Error chain: %v\n", err)

		// Check if it matches our custom sentinel
		if errors.Is(err, ErrConfigMissing) {
			fmt.Println("Matched ErrConfigMissing!")
		}

		// Check if it also matches standard library os.ErrNotExist
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("Matched standard os.ErrNotExist!")
		}
	}
}
