package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrInsufficientFunds = errors.New("insufficient account funds")

// TransferAtomic performs an atomic transfer between two accounts inside a single transaction.
func TransferAtomic(ctx context.Context, db *sql.DB, fromID, toID int, amount float64) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}
	defer tx.Rollback() // Safe deferred rollback

	// Check balance
	var currentBalance float64
	row := tx.QueryRowContext(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", fromID)
	if err := row.Scan(&currentBalance); err != nil {
		return fmt.Errorf("fetch balance failed: %w", err)
	}

	if currentBalance < amount {
		return ErrInsufficientFunds
	}

	// Debit sender
	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID); err != nil {
		return fmt.Errorf("debit failed: %w", err)
	}

	// Credit receiver
	if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID); err != nil {
		return fmt.Errorf("credit failed: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("ACID Transaction Transfer Example pattern initialized.")
}
