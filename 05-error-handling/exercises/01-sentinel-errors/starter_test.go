package userrepo

import (
	"errors"
	"testing"
)

func TestUserRepository(t *testing.T) {
	repo := NewRepository()

	// Test 1: Create user
	err := repo.CreateUser(1, "ashraf", "ashraf@domain.com")
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: CreateUser is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test 2: Duplicate email
	err = repo.CreateUser(2, "other", "ashraf@domain.com")
	if !errors.Is(err, ErrDuplicateEmail) {
		t.Errorf("expected ErrDuplicateEmail, got %v", err)
	}

	// Test 3: Find existing user
	u, err := repo.FindByID(1)
	if err != nil {
		t.Fatalf("unexpected error finding user: %v", err)
	}
	if u.Username != "ashraf" {
		t.Errorf("got username %s, want 'ashraf'", u.Username)
	}

	// Test 4: Find missing user with wrapped sentinel
	_, err = repo.FindByID(999)
	if !errors.Is(err, ErrUserNotFound) {
		t.Errorf("expected wrapped ErrUserNotFound, got %v", err)
	}
}
