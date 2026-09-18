package transfer

import (
	"context"
	"errors"
	"testing"
)

type mockTx struct {
	balances   map[string]float64
	committed  bool
	rolledBack bool
}

func (m *mockTx) GetBalance(ctx context.Context, accountID string) (float64, error) {
	b, ok := m.balances[accountID]
	if !ok {
		return 0, errors.New("account not found")
	}
	return b, nil
}

func (m *mockTx) UpdateBalance(ctx context.Context, accountID string, delta float64) error {
	m.balances[accountID] += delta
	return nil
}

func (m *mockTx) Commit() error {
	m.committed = true
	return nil
}

func (m *mockTx) Rollback() error {
	m.rolledBack = true
	return nil
}

func TestAtomicTransfer(t *testing.T) {
	ctx := context.Background()

	// 1. Check unimplemented starter
	testTx := &mockTx{balances: map[string]float64{"A": 100, "B": 50}}
	err := AtomicTransfer(ctx, testTx, "A", "B", 20)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: AtomicTransfer is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !testTx.committed {
		t.Errorf("expected tx to be committed")
	}
	if testTx.balances["A"] != 80 || testTx.balances["B"] != 70 {
		t.Errorf("balances incorrect: %+v", testTx.balances)
	}

	// 2. Insufficient funds
	lowTx := &mockTx{balances: map[string]float64{"A": 10, "B": 50}}
	err = AtomicTransfer(ctx, lowTx, "A", "B", 100)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
	if !lowTx.rolledBack {
		t.Errorf("expected tx to be rolled back on insufficient funds")
	}

	// 3. Same account
	err = AtomicTransfer(ctx, testTx, "A", "A", 10)
	if !errors.Is(err, ErrSameAccount) {
		t.Errorf("expected ErrSameAccount, got %v", err)
	}

	// 4. Negative amount
	err = AtomicTransfer(ctx, testTx, "A", "B", -5)
	if !errors.Is(err, ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}
