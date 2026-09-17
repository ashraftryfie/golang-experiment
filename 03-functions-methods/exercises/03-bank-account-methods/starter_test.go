package bank

import (
	"errors"
	"testing"
)

func TestBankAccount(t *testing.T) {
	acc, err := NewAccount("Ashraf", 100.0)
	if err != nil {
		t.Fatalf("unexpected error creating account: %v", err)
	}

	// Deposit
	err = acc.Deposit(50.0)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: Deposit is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	if acc.Balance() != 150.0 {
		t.Errorf("balance = %v, want 150.0", acc.Balance())
	}

	// Withdraw
	err = acc.Withdraw(30.0)
	if err != nil {
		t.Fatalf("withdraw failed: %v", err)
	}
	if acc.Balance() != 120.0 {
		t.Errorf("balance = %v, want 120.0", acc.Balance())
	}

	// Overdraw error
	err = acc.Withdraw(500.0)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}

	// Statement slice encapsulation check
	stmt := acc.Statement()
	if len(stmt) != 2 {
		t.Errorf("expected 2 transactions, got %d", len(stmt))
	}
}

func TestNewAccountValidations(t *testing.T) {
	if _, err := NewAccount("", 100); !errors.Is(err, ErrInvalidOwner) {
		t.Errorf("expected ErrInvalidOwner, got %v", err)
	}
	if _, err := NewAccount("Ashraf", -10); !errors.Is(err, ErrNegativeAmount) {
		t.Errorf("expected ErrNegativeAmount, got %v", err)
	}
}
