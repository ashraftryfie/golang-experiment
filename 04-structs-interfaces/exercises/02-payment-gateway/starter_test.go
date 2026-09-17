package payment

import (
	"errors"
	"testing"
)

func TestPaymentGateway(t *testing.T) {
	gw := NewPaymentGateway()
	gw.Register("stripe", CreditCardProcessor{CardNumber: "1234567812345678"})
	gw.Register("paypal", PayPalProcessor{Email: "user@example.com"})

	res, err := gw.Pay("stripe", 99.99)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: Pay is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("stripe payment failed: %v", err)
	}
	if !res.Success {
		t.Errorf("expected payment to succeed")
	}

	// Unsupported provider check
	_, err = gw.Pay("bitcoin", 50.0)
	if !errors.Is(err, ErrUnsupportedMethod) {
		t.Errorf("expected ErrUnsupportedMethod, got %v", err)
	}
}
