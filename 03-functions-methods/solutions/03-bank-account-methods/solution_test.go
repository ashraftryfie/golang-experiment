package bank

import (
	"errors"
	"testing"
)

func TestSolutionBankAccount(t *testing.T) {
	acc, err := NewAccount("Alice", 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := acc.Deposit(50); err != nil {
		t.Fatalf("deposit error: %v", err)
	}
	if err := acc.Withdraw(30); err != nil {
		t.Fatalf("withdraw error: %v", err)
	}
	if acc.Balance() != 120 {
		t.Errorf("balance = %v, want 120", acc.Balance())
	}

	// Negative amount check
	if err := acc.Deposit(-5); !errors.Is(err, ErrNegativeAmount) {
		t.Errorf("expected ErrNegativeAmount, got %v", err)
	}

	// Statement isolation
	stmt := acc.Statement()
	stmt[0].Amount = 99999
	if acc.Statement()[0].Amount == 99999 {
		t.Errorf("statement returned direct reference instead of copy!")
	}
}
