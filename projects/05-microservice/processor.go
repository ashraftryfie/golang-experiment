package orderprocessor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrDuplicateRequest = errors.New("duplicate request ignored (idempotent)")
	ErrOrderNotFound    = errors.New("order not found")
	ErrInvalidStatus    = errors.New("invalid order status transition")
)

type ProcessState string

const (
	StatePending   ProcessState = "PENDING"
	StateConfirmed ProcessState = "CONFIRMED"
	StateFulfilled ProcessState = "FULFILLED"
)

type OrderRecord struct {
	ID             string       `json:"id"`
	IdempotencyKey string       `json:"idempotency_key"`
	AmountCents    int64        `json:"amount_cents"`
	State          ProcessState `json:"state"`
	CreatedAt      time.Time    `json:"created_at"`
}

type OrderProcessor struct {
	mu          sync.RWMutex
	orders      map[string]*OrderRecord
	idempotency map[string]string // key -> orderID
}

func NewOrderProcessor() *OrderProcessor {
	return &OrderProcessor{
		orders:      make(map[string]*OrderRecord),
		idempotency: make(map[string]string),
	}
}

func (p *OrderProcessor) ProcessOrder(_ context.Context, orderID, idempotencyKey string, amountCents int64) (*OrderRecord, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check idempotency
	if existingID, exists := p.idempotency[idempotencyKey]; exists {
		return p.orders[existingID], ErrDuplicateRequest
	}

	record := &OrderRecord{
		ID:             orderID,
		IdempotencyKey: idempotencyKey,
		AmountCents:    amountCents,
		State:          StateConfirmed,
		CreatedAt:      time.Now().UTC(),
	}

	p.orders[orderID] = record
	p.idempotency[idempotencyKey] = orderID

	return record, nil
}

func (p *OrderProcessor) FulfillOrder(_ context.Context, orderID string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	order, exists := p.orders[orderID]
	if !exists {
		return ErrOrderNotFound
	}

	if order.State != StateConfirmed {
		return fmt.Errorf("%w: cannot fulfill from %s", ErrInvalidStatus, order.State)
	}

	order.State = StateFulfilled
	return nil
}

func (p *OrderProcessor) GetOrder(_ context.Context, orderID string) (*OrderRecord, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	order, exists := p.orders[orderID]
	if !exists {
		return nil, ErrOrderNotFound
	}
	clone := *order
	return &clone, nil
}
