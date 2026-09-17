package bank

import (
	"errors"
)

var (
	ErrInvalidOwner      = errors.New("owner name cannot be empty")
	ErrNegativeAmount    = errors.New("amount must be greater than zero")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrNotImplemented    = errors.New("TODO: implement method")
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
	return &BankAccount{
		owner:   owner,
		balance: initialDeposit,
	}, nil
}

func (a *BankAccount) Deposit(amount float64) error {
	// TODO: Validate amount > 0, update balance, append history
	return ErrNotImplemented
}

func (a *BankAccount) Withdraw(amount float64) error {
	// TODO: Validate amount > 0, check sufficient funds, update balance, append history
	return ErrNotImplemented
}

func (a BankAccount) Balance() float64 {
	return a.balance
}

func (a BankAccount) Statement() []Transaction {
	// TODO: Return a copy of history slice
	return nil
}
