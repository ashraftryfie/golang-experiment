package main

import (
	"errors"
	"testing"

	"github.com/ashraftryfie/golang-experiment/06-packages-modules/challenges/modular-monolith/internal/store"
	"github.com/ashraftryfie/golang-experiment/06-packages-modules/challenges/modular-monolith/pkg/validator"
)

func TestAppIntegration(t *testing.T) {
	st := store.NewMemoryStore()
	app := NewApp(st)

	// Valid creation
	err := app.CreateProduct("p1", "Desk Mat", 24.50)
	if err != nil {
		t.Fatalf("expected product creation to succeed, got %v", err)
	}

	// Retrieve
	p, err := st.Get("p1")
	if err != nil {
		t.Fatalf("expected retrieval to succeed, got %v", err)
	}
	if p.Name != "Desk Mat" {
		t.Errorf("got name %q, want 'Desk Mat'", p.Name)
	}

	// Validation rejection
	err = app.CreateProduct("p2", "", 10.0)
	if !errors.Is(err, validator.ErrEmptyName) {
		t.Errorf("expected ErrEmptyName, got %v", err)
	}
}
