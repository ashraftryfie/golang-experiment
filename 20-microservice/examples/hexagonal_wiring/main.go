package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ==========================================
// 1. DOMAIN LAYER (Zero external dependencies)
// ==========================================

var (
	ErrNegativeAmount  = errors.New("amount must be positive")
	ErrInsufficientBal = errors.New("insufficient balance")
	ErrUserNotFound    = errors.New("user not found")
)

// Money is an immutable Value Object (stored in cents)
type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	if cents < 0 {
		return Money{}, ErrNegativeAmount
	}
	return Money{cents: cents, currency: currency}, nil
}

func (m Money) Cents() int64     { return m.cents }
func (m Money) Currency() string { return m.currency }
func (m Money) String() string   { return fmt.Sprintf("%.2f %s", float64(m.cents)/100.0, m.currency) }

// UserAccount is an Entity with distinct identity and domain invariants
type UserAccount struct {
	id      string
	name    string
	balance Money
}

func NewUserAccount(id, name string, initialBalance Money) *UserAccount {
	return &UserAccount{
		id:      id,
		name:    name,
		balance: initialBalance,
	}
}

func (u *UserAccount) ID() string     { return u.id }
func (u *UserAccount) Name() string   { return u.name }
func (u *UserAccount) Balance() Money { return u.balance }

func (u *UserAccount) Deposit(m Money) {
	u.balance.cents += m.cents
}

func (u *UserAccount) Withdraw(m Money) error {
	if u.balance.cents < m.cents {
		return ErrInsufficientBal
	}
	u.balance.cents -= m.cents
	return nil
}

// ==========================================
// 2. PORTS (Interfaces defined by consumer)
// ==========================================

// AccountRepository is the primary persistence port
type AccountRepository interface {
	FindByID(ctx context.Context, id string) (*UserAccount, error)
	Save(ctx context.Context, account *UserAccount) error
}

// AuditLogger is a secondary port for compliance logging
type AuditLogger interface {
	LogEvent(ctx context.Context, action string, details map[string]any)
}

// ==========================================
// 3. SERVICE LAYER (Use-case orchestration)
// ==========================================

type BankingService struct {
	repo   AccountRepository
	logger AuditLogger
}

// NewBankingService uses constructor-based Dependency Injection
func NewBankingService(repo AccountRepository, logger AuditLogger) *BankingService {
	return &BankingService{
		repo:   repo,
		logger: logger,
	}
}

func (s *BankingService) Transfer(ctx context.Context, fromID, toID string, amount Money) error {
	from, err := s.repo.FindByID(ctx, fromID)
	if err != nil {
		return fmt.Errorf("sender lookup failed: %w", err)
	}

	to, err := s.repo.FindByID(ctx, toID)
	if err != nil {
		return fmt.Errorf("receiver lookup failed: %w", err)
	}

	if err := from.Withdraw(amount); err != nil {
		return err
	}
	to.Deposit(amount)

	if err := s.repo.Save(ctx, from); err != nil {
		return fmt.Errorf("failed saving sender: %w", err)
	}
	if err := s.repo.Save(ctx, to); err != nil {
		return fmt.Errorf("failed saving receiver: %w", err)
	}

	s.logger.LogEvent(ctx, "transfer.completed", map[string]any{
		"from_id": fromID,
		"to_id":   toID,
		"amount":  amount.String(),
	})

	return nil
}

// ==========================================
// 4. ADAPTERS (Infrastructure implementations)
// ==========================================

// InMemoryAccountRepo is an in-memory adapter implementing AccountRepository
type InMemoryAccountRepo struct {
	mu       sync.RWMutex
	accounts map[string]*UserAccount
}

func NewInMemoryAccountRepo() *InMemoryAccountRepo {
	return &InMemoryAccountRepo{accounts: make(map[string]*UserAccount)}
}

func (r *InMemoryAccountRepo) FindByID(_ context.Context, id string) (*UserAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	acc, exists := r.accounts[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	// Return a copy to preserve repository boundary integrity
	clone := *acc
	return &clone, nil
}

func (r *InMemoryAccountRepo) Save(_ context.Context, account *UserAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	clone := *account
	r.accounts[account.ID()] = &clone
	return nil
}

// StdoutAuditLogger is an adapter implementing AuditLogger
type StdoutAuditLogger struct{}

func (l *StdoutAuditLogger) LogEvent(_ context.Context, action string, details map[string]any) {
	fmt.Printf("[AUDIT %s] action=%s details=%v\n", time.Now().Format(time.RFC3339), action, details)
}

// ==========================================
// 5. APPLICATION COMPOSITION ROOT (Wiring)
// ==========================================

func main() {
	ctx := context.Background()

	// 1. Initialize Adapters
	repo := NewInMemoryAccountRepo()
	logger := &StdoutAuditLogger{}

	// Seed domain data
	usd100, _ := NewMoney(10000, "USD")
	usd50, _ := NewMoney(5000, "USD")
	_ = repo.Save(ctx, NewUserAccount("acc-1", "Alice", usd100))
	_ = repo.Save(ctx, NewUserAccount("acc-2", "Bob", usd50))

	// 2. Inject Adapters into Service
	svc := NewBankingService(repo, logger)

	// 3. Execute Use Case
	transferAmount, _ := NewMoney(2500, "USD")
	fmt.Println("Executing transfer of", transferAmount, "from Alice to Bob...")
	if err := svc.Transfer(ctx, "acc-1", "acc-2", transferAmount); err != nil {
		fmt.Printf("Transfer error: %v\n", err)
		return
	}

	alice, _ := repo.FindByID(ctx, "acc-1")
	bob, _ := repo.FindByID(ctx, "acc-2")
	fmt.Printf("Alice's new balance: %s\n", alice.Balance())
	fmt.Printf("Bob's new balance:   %s\n", bob.Balance())
}
