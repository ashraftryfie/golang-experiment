package payment

import (
	"errors"
)

var (
	ErrInvalidCard       = errors.New("invalid card number (must be 16 digits)")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrInvalidWallet     = errors.New("invalid crypto wallet address")
	ErrUnsupportedMethod = errors.New("unsupported payment method")
	ErrNotImplemented    = errors.New("TODO: implement")
)

type PaymentResult struct {
	TransactionID string
	Success       bool
	Message       string
}

type PaymentProcessor interface {
	ProcessPayment(amount float64) (PaymentResult, error)
}

type CreditCardProcessor struct {
	CardNumber string
}

func (c CreditCardProcessor) ProcessPayment(amount float64) (PaymentResult, error) {
	// TODO: Validate card and return result
	return PaymentResult{}, ErrNotImplemented
}

type PayPalProcessor struct {
	Email string
}

func (p PayPalProcessor) ProcessPayment(amount float64) (PaymentResult, error) {
	// TODO: Validate email and return result
	return PaymentResult{}, ErrNotImplemented
}

type PaymentGateway struct {
	processors map[string]PaymentProcessor
}

func NewPaymentGateway() *PaymentGateway {
	return &PaymentGateway{
		processors: make(map[string]PaymentProcessor),
	}
}

func (g *PaymentGateway) Register(name string, p PaymentProcessor) {
	g.processors[name] = p
}

func (g *PaymentGateway) Pay(provider string, amount float64) (PaymentResult, error) {
	// TODO: Lookup provider and process payment
	return PaymentResult{}, ErrNotImplemented
}
