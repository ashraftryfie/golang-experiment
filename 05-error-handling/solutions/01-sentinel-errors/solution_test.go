package userrepo

import (
	"errors"
	"testing"
)

func TestSolutionUserRepository(t *testing.T) {
	repo := NewRepository()

	if err := repo.CreateUser(10, "gopher", "gopher@golang.org"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	u, err := repo.FindByID(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.Username != "gopher" {
		t.Errorf("got %s, want 'gopher'", u.Username)
	}

	_, err = repo.FindByID(99)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
