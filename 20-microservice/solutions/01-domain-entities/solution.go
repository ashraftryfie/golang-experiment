package domainentities

import (
	"errors"
	"fmt"
)

var (
	ErrNegativeAmount      = errors.New("negative amount not allowed")
	ErrInvalidCurrency     = errors.New("currency code cannot be empty")
	ErrCurrencyMismatch    = errors.New("currency mismatch")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidID           = errors.New("id cannot be empty")
)

// Money is an immutable Value Object
type Money struct {
	cents    int64
	currency string
}

func NewMoney(cents int64, currency string) (Money, error) {
	if cents < 0 {
		return Money{}, ErrNegativeAmount
	}
	if currency == "" {
		return Money{}, ErrInvalidCurrency
	}
	return Money{cents: cents, currency: currency}, nil
}

func (m Money) Cents() int64 {
	return m.cents
}

func (m Money) Currency() string {
	return m.currency
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{cents: m.cents + other.cents, currency: m.currency}, nil
}

func (m Money) Subtract(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	if m.cents < other.cents {
		return Money{}, ErrInsufficientBalance
	}
	return Money{cents: m.cents - other.cents, currency: m.currency}, nil
}

func (m Money) Equals(other Money) bool {
	return m.cents == other.cents && m.currency == other.currency
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.cents)/100.0, m.currency)
}

// Wallet is an Aggregate Entity
type Wallet struct {
	id      string
	ownerID string
	balance Money
}

func NewWallet(id, ownerID string, initialBalance Money) (*Wallet, error) {
	if id == "" || ownerID == "" {
		return nil, ErrInvalidID
	}
	return &Wallet{
		id:      id,
		ownerID: ownerID,
		balance: initialBalance,
	}, nil
}

func (w *Wallet) ID() string {
	return w.id
}

func (w *Wallet) OwnerID() string {
	return w.ownerID
}

func (w *Wallet) Balance() Money {
	return w.balance
}

func (w *Wallet) Deposit(m Money) error {
	newBal, err := w.balance.Add(m)
	if err != nil {
		return err
	}
	w.balance = newBal
	return nil
}

func (w *Wallet) Withdraw(m Money) error {
	newBal, err := w.balance.Subtract(m)
	if err != nil {
		return err
	}
	w.balance = newBal
	return nil
}
