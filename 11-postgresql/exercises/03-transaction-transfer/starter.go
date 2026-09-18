package transfer

import (
	"context"
	"errors"
)

var (
	ErrNotImplemented    = errors.New("TODO: implement AtomicTransfer")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrSameAccount       = errors.New("cannot transfer to self")
	ErrInsufficientFunds = errors.New("insufficient account funds")
)

type TxExecutor interface {
	GetBalance(ctx context.Context, accountID string) (float64, error)
	UpdateBalance(ctx context.Context, accountID string, delta float64) error
	Commit() error
	Rollback() error
}

// AtomicTransfer orchestrates a safe transactional transfer between accounts.
func AtomicTransfer(ctx context.Context, tx TxExecutor, fromID, toID string, amount float64) error {
	// TODO: Validate, defer tx.Rollback(), verify balance, debit, credit, tx.Commit()
	return ErrNotImplemented
}
