package transfer_solution

import (
	"context"
	"errors"
	"fmt"
)

var (
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
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if fromID == toID {
		return ErrSameAccount
	}

	defer tx.Rollback() // Safe deferred rollback

	currentBalance, err := tx.GetBalance(ctx, fromID)
	if err != nil {
		return fmt.Errorf("fetch sender balance failed: %w", err)
	}

	if currentBalance < amount {
		return ErrInsufficientFunds
	}

	if err := tx.UpdateBalance(ctx, fromID, -amount); err != nil {
		return fmt.Errorf("debit failed: %w", err)
	}

	if err := tx.UpdateBalance(ctx, toID, amount); err != nil {
		return fmt.Errorf("credit failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}
