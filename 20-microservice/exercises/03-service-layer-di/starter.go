package servicelayer

import (
	"context"
	"errors"
)

var (
	ErrInvalidAmount     = errors.New("transfer amount must be positive")
	ErrSameAccount       = errors.New("cannot transfer to the same account")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
	ErrNotImplemented    = errors.New("not implemented yet")
)

type Account struct {
	id      string
	balance int64
}

func NewAccount(id string, balance int64) *Account {
	return &Account{id: id, balance: balance}
}

func (a *Account) ID() string      { return a.id }
func (a *Account) Balance() int64  { return a.balance }

func (a *Account) Debit(amount int64) error {
	// TODO: Check balance and debit
	return ErrNotImplemented
}

func (a *Account) Credit(amount int64) {
	// TODO: Credit amount
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
	// TODO: Validate, load accounts, transfer, save with rollback compensation, notify
	return ErrNotImplemented
}
