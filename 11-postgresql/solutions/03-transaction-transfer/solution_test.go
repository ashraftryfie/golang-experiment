package transfer_solution

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

func TestAtomicTransferSolution(t *testing.T) {
	ctx := context.Background()
	testTx := &mockTx{balances: map[string]float64{"A": 100, "B": 50}}

	err := AtomicTransfer(ctx, testTx, "A", "B", 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !testTx.committed {
		t.Errorf("expected commit")
	}
	if testTx.balances["A"] != 80 || testTx.balances["B"] != 70 {
		t.Errorf("balances incorrect: %+v", testTx.balances)
	}

	// Insufficient funds
	lowTx := &mockTx{balances: map[string]float64{"A": 10, "B": 50}}
	err = AtomicTransfer(ctx, lowTx, "A", "B", 100)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}
