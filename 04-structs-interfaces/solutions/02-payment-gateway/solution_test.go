package payment

import (
	"errors"
	"testing"
)

func TestSolutionPaymentGateway(t *testing.T) {
	gw := NewPaymentGateway()
	gw.Register("stripe", CreditCardProcessor{CardNumber: "1111222233334444"})
	gw.Register("paypal", PayPalProcessor{Email: "ashraf@test.com"})

	res, err := gw.Pay("stripe", 50.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success {
		t.Errorf("expected success")
	}

	res2, err := gw.Pay("paypal", 25.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res2.Success {
		t.Errorf("expected success")
	}

	// Bad card
	gw.Register("bad_card", CreditCardProcessor{CardNumber: "123"})
	_, err = gw.Pay("bad_card", 10.0)
	if !errors.Is(err, ErrInvalidCard) {
		t.Errorf("expected ErrInvalidCard, got %v", err)
	}
}
