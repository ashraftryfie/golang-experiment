package bank

import (
	"errors"
)

var (
	ErrInvalidOwner      = errors.New("owner name cannot be empty")
	ErrNegativeAmount    = errors.New("amount must be greater than zero")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type TransactionType string

const (
	TxDeposit  TransactionType = "DEPOSIT"
	TxWithdraw TransactionType = "WITHDRAW"
)

type Transaction struct {
	Type   TransactionType
	Amount float64
}

type BankAccount struct {
	owner   string
	balance float64
	history []Transaction
}

func NewAccount(owner string, initialDeposit float64) (*BankAccount, error) {
	if owner == "" {
		return nil, ErrInvalidOwner
	}
	if initialDeposit < 0 {
		return nil, ErrNegativeAmount
	}

	acc := &BankAccount{
		owner:   owner,
		balance: initialDeposit,
		history: make([]Transaction, 0),
	}

	if initialDeposit > 0 {
		acc.history = append(acc.history, Transaction{
			Type:   TxDeposit,
			Amount: initialDeposit,
		})
	}

	return acc, nil
}

func (a *BankAccount) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrNegativeAmount
	}
	a.balance += amount
	a.history = append(a.history, Transaction{
		Type:   TxDeposit,
		Amount: amount,
	})
	return nil
}

func (a *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrNegativeAmount
	}
	if amount > a.balance {
		return ErrInsufficientFunds
	}
	a.balance -= amount
	a.history = append(a.history, Transaction{
		Type:   TxWithdraw,
		Amount: amount,
	})
	return nil
}

func (a BankAccount) Balance() float64 {
	return a.balance
}

func (a BankAccount) Statement() []Transaction {
	clone := make([]Transaction, len(a.history))
	copy(clone, a.history)
	return clone
}
