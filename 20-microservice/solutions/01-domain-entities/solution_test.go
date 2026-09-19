package domainentities

import (
	"errors"
	"testing"
)

func TestMoney_CreationAndValidation(t *testing.T) {
	tests := []struct {
		name      string
		cents     int64
		currency  string
		expectErr error
	}{
		{"valid money", 1500, "USD", nil},
		{"zero amount is valid", 0, "EUR", nil},
		{"negative amount rejected", -50, "USD", ErrNegativeAmount},
		{"empty currency rejected", 100, "", ErrInvalidCurrency},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m, err := NewMoney(tc.cents, tc.currency)
			if tc.expectErr != nil {
				if !errors.Is(err, tc.expectErr) {
					t.Fatalf("expected error %v, got %v", tc.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.Cents() != tc.cents || m.Currency() != tc.currency {
				t.Fatalf("unexpected values: cents=%d, currency=%s", m.Cents(), m.Currency())
			}
		})
	}
}

func TestMoney_Operations(t *testing.T) {
	m1, _ := NewMoney(1050, "USD")
	m2, _ := NewMoney(450, "USD")
	eur, _ := NewMoney(500, "EUR")

	// Add
	sum, err := m1.Add(m2)
	if err != nil || sum.Cents() != 1500 {
		t.Fatalf("Add failed: sum=%v, err=%v", sum, err)
	}

	// Add mismatch
	_, err = m1.Add(eur)
	if !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("expected ErrCurrencyMismatch, got %v", err)
	}

	// Subtract
	diff, err := m1.Subtract(m2)
	if err != nil || diff.Cents() != 600 {
		t.Fatalf("Subtract failed: diff=%v, err=%v", diff, err)
	}

	// Subtract overdraft
	_, err = m2.Subtract(m1)
	if !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}

	// Equals
	m1Copy, _ := NewMoney(1050, "USD")
	if !m1.Equals(m1Copy) {
		t.Fatal("expected m1 to equal m1Copy")
	}
	if m1.Equals(m2) {
		t.Fatal("expected m1 not to equal m2")
	}

	// String format
	if m1.String() != "10.50 USD" {
		t.Fatalf("expected '10.50 USD', got '%s'", m1.String())
	}
}

func TestWallet_Invariants(t *testing.T) {
	bal, _ := NewMoney(5000, "USD")
	wallet, err := NewWallet("w-1", "user-1", bal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Deposit
	dep, _ := NewMoney(2000, "USD")
	if err := wallet.Deposit(dep); err != nil {
		t.Fatalf("deposit failed: %v", err)
	}
	if wallet.Balance().Cents() != 7000 {
		t.Fatalf("expected 7000 cents, got %d", wallet.Balance().Cents())
	}

	// Withdraw
	wit, _ := NewMoney(3000, "USD")
	if err := wallet.Withdraw(wit); err != nil {
		t.Fatalf("withdraw failed: %v", err)
	}
	if wallet.Balance().Cents() != 4000 {
		t.Fatalf("expected 4000 cents, got %d", wallet.Balance().Cents())
	}

	// Withdraw over balance
	bigWit, _ := NewMoney(10000, "USD")
	if err := wallet.Withdraw(bigWit); !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}

	// Deposit wrong currency
	depEur, _ := NewMoney(500, "EUR")
	if err := wallet.Deposit(depEur); !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("expected ErrCurrencyMismatch, got %v", err)
	}
}
