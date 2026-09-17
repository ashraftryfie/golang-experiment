package payment

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidCard       = errors.New("invalid card number (must be 16 digits)")
	ErrInvalidEmail      = errors.New("invalid email address")
	ErrUnsupportedMethod = errors.New("unsupported payment method")
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
	if len(c.CardNumber) != 16 {
		return PaymentResult{}, ErrInvalidCard
	}
	return PaymentResult{
		TransactionID: fmt.Sprintf("cc_tx_%d", int(amount*100)),
		Success:       true,
		Message:       fmt.Sprintf("Charged $%.2f to Card Ending in %s", amount, c.CardNumber[12:]),
	}, nil
}

type PayPalProcessor struct {
	Email string
}

func (p PayPalProcessor) ProcessPayment(amount float64) (PaymentResult, error) {
	if !strings.Contains(p.Email, "@") {
		return PaymentResult{}, ErrInvalidEmail
	}
	return PaymentResult{
		TransactionID: fmt.Sprintf("pp_tx_%d", int(amount*100)),
		Success:       true,
		Message:       fmt.Sprintf("Processed $%.2f via PayPal to %s", amount, p.Email),
	}, nil
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
	proc, exists := g.processors[provider]
	if !exists {
		return PaymentResult{}, ErrUnsupportedMethod
	}
	return proc.ProcessPayment(amount)
}
