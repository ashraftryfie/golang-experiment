package orderservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrder_StateTransitions(t *testing.T) {
	items := []OrderItem{
		{ProductID: "prod-1", Quantity: 2, UnitPriceCents: 1500},
		{ProductID: "prod-2", Quantity: 1, UnitPriceCents: 2000},
	}

	order, err := NewOrder("ord-1", "cust-1", items)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Invariant: TotalCents = 2*1500 + 1*2000 = 5000 cents ($50.00)
	if order.TotalCents() != 5000 {
		t.Fatalf("expected 5000 cents, got %d", order.TotalCents())
	}

	if order.Status != StatusPending {
		t.Fatalf("expected initial status PENDING, got %s", order.Status)
	}

	// State transition: Pending -> Paid
	if err := order.MarkPaid(); err != nil {
		t.Fatalf("MarkPaid failed: %v", err)
	}
	if order.Status != StatusPaid {
		t.Fatalf("expected status PAID, got %s", order.Status)
	}

	// Cannot pay twice
	if err := order.MarkPaid(); !errors.Is(err, ErrAlreadyPaid) {
		t.Fatalf("expected ErrAlreadyPaid, got %v", err)
	}

	// State transition: Paid -> Shipped
	if err := order.MarkShipped(); err != nil {
		t.Fatalf("MarkShipped failed: %v", err)
	}
	if order.Status != StatusShipped {
		t.Fatalf("expected status SHIPPED, got %s", order.Status)
	}

	// Cannot cancel after shipping
	if err := order.Cancel(); !errors.Is(err, ErrCannotCancelShipped) {
		t.Fatalf("expected ErrCannotCancelShipped, got %v", err)
	}
}

func TestOrderService_EndToEndLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryOrderRepository()
	payment := NewMockPaymentGateway()
	dispatcher := NewMockEventDispatcher()

	svc := NewOrderService(repo, payment, dispatcher)
	svc.SetIDGenerator(func() string { return "order-fixed-101" })

	// 1. Create order
	items := []OrderItem{
		{ProductID: "keyboard", Quantity: 1, UnitPriceCents: 8500},
	}
	order, err := svc.CreateOrder(ctx, "cust-alice", items)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}
	if order.ID != "order-fixed-101" || order.Status != StatusPending {
		t.Fatalf("unexpected order after creation: %+v", order)
	}

	// 2. Pay order
	if err := svc.PayOrder(ctx, order.ID); err != nil {
		t.Fatalf("PayOrder failed: %v", err)
	}

	paidOrder, _ := repo.GetByID(ctx, order.ID)
	if paidOrder.Status != StatusPaid {
		t.Fatalf("expected status PAID, got %s", paidOrder.Status)
	}
	if len(payment.Charges) != 1 || payment.Charges[0] != 8500 {
		t.Fatalf("unexpected payment charge: %v", payment.Charges)
	}

	// 3. Ship order
	if err := svc.ShipOrder(ctx, order.ID); err != nil {
		t.Fatalf("ShipOrder failed: %v", err)
	}

	// 4. Verify dispatched events
	expectedEvents := []string{"order.created", "order.paid", "order.shipped"}
	if len(dispatcher.Events) != 3 {
		t.Fatalf("expected 3 events, got %v", dispatcher.Events)
	}
	for i, exp := range expectedEvents {
		if dispatcher.Events[i] != exp {
			t.Errorf("expected event %d to be %s, got %s", i, exp, dispatcher.Events[i])
		}
	}
}

func TestOrderService_PaymentFailure(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryOrderRepository()
	payment := NewMockPaymentGateway()
	payment.ShouldFail = true // Payment gateway returns card declined
	dispatcher := NewMockEventDispatcher()

	svc := NewOrderService(repo, payment, dispatcher)
	order, _ := svc.CreateOrder(ctx, "cust-bob", []OrderItem{
		{ProductID: "mouse", Quantity: 1, UnitPriceCents: 2500},
	})

	err := svc.PayOrder(ctx, order.ID)
	if !errors.Is(err, ErrPaymentFailed) {
		t.Fatalf("expected ErrPaymentFailed, got %v", err)
	}

	// Order must remain PENDING because payment failed
	stored, _ := repo.GetByID(ctx, order.ID)
	if stored.Status != StatusPending {
		t.Fatalf("order status should remain PENDING after failed payment, got %s", stored.Status)
	}
}

func TestOrderHTTPHandler_Integration(t *testing.T) {
	repo := NewInMemoryOrderRepository()
	payment := NewMockPaymentGateway()
	dispatcher := NewMockEventDispatcher()

	svc := NewOrderService(repo, payment, dispatcher)
	svc.SetIDGenerator(func() string { return "order-http-999" })
	handler := NewOrderHTTPHandler(svc)

	server := httptest.NewServer(handler)
	defer server.Close()

	client := server.Client()

	// 1. POST /orders
	reqBody := `{"customer_id":"cust-charlie","items":[{"product_id":"monitor","quantity":1,"unit_price_cents":25000}]}`
	resp, err := client.Post(server.URL+"/orders", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("POST /orders request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	var created Order
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("failed decoding created order: %v", err)
	}
	resp.Body.Close()

	if created.ID != "order-http-999" || created.TotalCents() != 25000 {
		t.Fatalf("unexpected created order: %+v", created)
	}

	// 2. POST /orders/order-http-999/pay
	payResp, err := client.Post(server.URL+"/orders/order-http-999/pay", "application/json", nil)
	if err != nil {
		t.Fatalf("POST pay failed: %v", err)
	}
	if payResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", payResp.StatusCode)
	}
	payResp.Body.Close()

	// 3. GET /orders/order-http-999
	getResp, err := client.Get(server.URL + "/orders/order-http-999")
	if err != nil {
		t.Fatalf("GET /orders request failed: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", getResp.StatusCode)
	}

	var fetched Order
	_ = json.NewDecoder(getResp.Body).Decode(&fetched)
	getResp.Body.Close()

	if fetched.Status != StatusPaid {
		t.Fatalf("expected fetched order to be PAID, got %s", fetched.Status)
	}
}
