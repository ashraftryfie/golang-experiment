package servicelayer

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type mockRepo struct {
	mu       sync.Mutex
	accounts map[string]*Account
	failOn   string // ID that triggers update failure
}

func newMockRepo() *mockRepo {
	return &mockRepo{accounts: make(map[string]*Account)}
}

func (m *mockRepo) Get(_ context.Context, id string) (*Account, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	acc, ok := m.accounts[id]
	if !ok {
		return nil, ErrAccountNotFound
	}
	clone := *acc
	return &clone, nil
}

func (m *mockRepo) Update(_ context.Context, account *Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.failOn == account.ID() {
		return errors.New("simulated database failure")
	}
	clone := *account
	m.accounts[account.ID()] = &clone
	return nil
}

type mockNotifier struct {
	mu            sync.Mutex
	notifications []string
}

func (n *mockNotifier) NotifyTransfer(_ context.Context, toID string, _ int64) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.notifications = append(n.notifications, toID)
	return nil
}

func TestTransferService_Success(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	repo.accounts["acc-1"] = NewAccount("acc-1", 1000)
	repo.accounts["acc-2"] = NewAccount("acc-2", 500)

	notifier := &mockNotifier{}
	svc := NewTransferService(repo, notifier)

	err := svc.ExecuteTransfer(ctx, "acc-1", "acc-2", 300)
	if err != nil {
		t.Fatalf("unexpected transfer error: %v", err)
	}

	acc1, _ := repo.Get(ctx, "acc-1")
	acc2, _ := repo.Get(ctx, "acc-2")

	if acc1.Balance() != 700 {
		t.Fatalf("expected acc-1 balance 700, got %d", acc1.Balance())
	}
	if acc2.Balance() != 800 {
		t.Fatalf("expected acc-2 balance 800, got %d", acc2.Balance())
	}

	if len(notifier.notifications) != 1 || notifier.notifications[0] != "acc-2" {
		t.Fatalf("expected notification to acc-2, got %v", notifier.notifications)
	}
}

func TestTransferService_InsufficientFunds(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	repo.accounts["acc-1"] = NewAccount("acc-1", 100)
	repo.accounts["acc-2"] = NewAccount("acc-2", 500)

	svc := NewTransferService(repo, nil)

	err := svc.ExecuteTransfer(ctx, "acc-1", "acc-2", 300)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestTransferService_SameAccount(t *testing.T) {
	ctx := context.Background()
	svc := NewTransferService(nil, nil)

	err := svc.ExecuteTransfer(ctx, "acc-1", "acc-1", 100)
	if !errors.Is(err, ErrSameAccount) {
		t.Fatalf("expected ErrSameAccount, got %v", err)
	}
}

func TestTransferService_CompensatingRollback(t *testing.T) {
	ctx := context.Background()
	repo := newMockRepo()
	repo.accounts["acc-1"] = NewAccount("acc-1", 1000)
	repo.accounts["acc-2"] = NewAccount("acc-2", 500)
	repo.failOn = "acc-2" // Simulate database failure when updating receiver

	svc := NewTransferService(repo, nil)

	err := svc.ExecuteTransfer(ctx, "acc-1", "acc-2", 400)
	if err == nil {
		t.Fatal("expected error on receiver update failure")
	}

	// Verify sender balance was rolled back to original 1000
	acc1, _ := repo.Get(ctx, "acc-1")
	if acc1.Balance() != 1000 {
		t.Fatalf("expected sender balance to be restored to 1000, got %d", acc1.Balance())
	}
}
