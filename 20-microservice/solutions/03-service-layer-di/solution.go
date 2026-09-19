package servicelayer

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount     = errors.New("transfer amount must be positive")
	ErrSameAccount       = errors.New("cannot transfer to the same account")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
)

type Account struct {
	id      string
	balance int64
}

func NewAccount(id string, balance int64) *Account {
	return &Account{id: id, balance: balance}
}

func (a *Account) ID() string     { return a.id }
func (a *Account) Balance() int64 { return a.balance }

func (a *Account) Debit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.balance < amount {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	return nil
}

func (a *Account) Credit(amount int64) {
	a.balance += amount
}

type AccountRepository interface {
	Get(ctx context.Context, id string) (*Account, error)
	Update(ctx context.Context, account *Account) error
}

type Notifier interface {
	NotifyTransfer(ctx context.Context, toID string, amount int64) error
}

type TransferService struct {
	repo     AccountRepository
	notifier Notifier
}

func NewTransferService(repo AccountRepository, notifier Notifier) *TransferService {
	return &TransferService{
		repo:     repo,
		notifier: notifier,
	}
}

func (s *TransferService) ExecuteTransfer(ctx context.Context, fromID, toID string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if fromID == toID {
		return ErrSameAccount
	}

	fromAcc, err := s.repo.Get(ctx, fromID)
	if err != nil {
		return fmt.Errorf("lookup sender %s: %w", fromID, err)
	}

	toAcc, err := s.repo.Get(ctx, toID)
	if err != nil {
		return fmt.Errorf("lookup receiver %s: %w", toID, err)
	}

	if err := fromAcc.Debit(amount); err != nil {
		return err
	}
	toAcc.Credit(amount)

	if err := s.repo.Update(ctx, fromAcc); err != nil {
		return fmt.Errorf("failed updating sender: %w", err)
	}

	if err := s.repo.Update(ctx, toAcc); err != nil {
		// Rollback sender debit as compensating transaction
		fromAcc.Credit(amount)
		_ = s.repo.Update(ctx, fromAcc)
		return fmt.Errorf("failed updating receiver (sender rolled back): %w", err)
	}

	if s.notifier != nil {
		_ = s.notifier.NotifyTransfer(ctx, toID, amount)
	}

	return nil
}
