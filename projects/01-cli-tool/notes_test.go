package main

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestNoteStore_CRUD(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-notes.json")

	store, err := NewNoteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create note store: %v", err)
	}

	// 1. Add Note
	n1, err := store.Add("Go Concurrency", "Goroutines and channels", []string{"go", "concurrency"})
	if err != nil {
		t.Fatalf("failed to add note: %v", err)
	}
	if n1.ID != 1 {
		t.Fatalf("expected ID 1, got %d", n1.ID)
	}

	// Add empty title fails
	_, err = store.Add("", "empty", nil)
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}

	// 2. Add second note
	_, _ = store.Add("Docker Basics", "Containers and compose", []string{"docker", "devops"})

	// 3. List
	notes := store.List()
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}

	// 4. Search
	searchRes := store.Search("goroutines")
	if len(searchRes) != 1 || searchRes[0].Title != "Go Concurrency" {
		t.Fatalf("expected 1 match for 'goroutines', got %v", searchRes)
	}

	tagRes := store.Search("devops")
	if len(tagRes) != 1 || tagRes[0].Title != "Docker Basics" {
		t.Fatalf("expected 1 match for tag 'devops', got %v", tagRes)
	}

	// 5. Delete
	if err := store.Delete(1); err != nil {
		t.Fatalf("failed to delete note 1: %v", err)
	}

	// Verify persistence reloaded
	reloaded, err := NewNoteStore(dbPath)
	if err != nil {
		t.Fatalf("failed reloading note store: %v", err)
	}
	if len(reloaded.List()) != 1 {
		t.Fatalf("expected 1 note after reload, got %d", len(reloaded.List()))
	}

	// Delete non-existent
	if err := reloaded.Delete(999); !errors.Is(err, ErrNoteNotFound) {
		t.Fatalf("expected ErrNoteNotFound, got %v", err)
	}
}
