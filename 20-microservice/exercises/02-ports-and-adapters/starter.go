package portsadapters

import (
	"context"
	"errors"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
	ErrInvalidSKU        = errors.New("sku cannot be empty")
	ErrNotImplemented    = errors.New("not implemented yet")
)

// Product is a domain entity
type Product struct {
	id    string
	sku   string
	stock int
}

func NewProduct(id, sku string, initialStock int) (*Product, error) {
	// TODO: Validate inputs and return Product
	return nil, ErrNotImplemented
}

func (p *Product) ID() string  { return p.id }
func (p *Product) SKU() string { return p.sku }
func (p *Product) Stock() int  { return p.stock }

func (p *Product) DeductStock(quantity int) error {
	// TODO: Enforce invariants and deduct
	return ErrNotImplemented
}

func (p *Product) AddStock(quantity int) error {
	// TODO: Enforce invariants and add
	return ErrNotImplemented
}

// ProductRepository is a persistence port
type ProductRepository interface {
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	Save(ctx context.Context, product *Product) error
}

// AlertPublisher is an outbound event port
type AlertPublisher interface {
	PublishLowStock(ctx context.Context, sku string, remaining int) error
}

// InMemoryProductRepo is an adapter implementing ProductRepository
type InMemoryProductRepo struct {
	// TODO: Add fields for thread-safe storage
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	// TODO: Initialize and return
	return nil
}

func (r *InMemoryProductRepo) GetBySKU(ctx context.Context, sku string) (*Product, error) {
	// TODO: Thread-safely retrieve cloned product
	return nil, ErrNotImplemented
}

func (r *InMemoryProductRepo) Save(ctx context.Context, product *Product) error {
	// TODO: Thread-safely store cloned product
	return ErrNotImplemented
}

// MockAlertPublisher is an adapter implementing AlertPublisher
type MockAlertPublisher struct {
	// TODO: Thread-safe slice or map to record alerts
}

func NewMockAlertPublisher() *MockAlertPublisher {
	// TODO: Initialize mock
	return nil
}

func (m *MockAlertPublisher) PublishLowStock(ctx context.Context, sku string, remaining int) error {
	// TODO: Record alert
	return ErrNotImplemented
}

func (m *MockAlertPublisher) RecordedAlerts() map[string]int {
	// TODO: Return copy of recorded alerts
	return nil
}
