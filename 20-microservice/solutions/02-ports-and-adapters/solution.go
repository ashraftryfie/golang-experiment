package portsadapters

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("quantity must be positive")
	ErrInvalidSKU        = errors.New("sku cannot be empty")
)

// Product is a domain entity
type Product struct {
	id    string
	sku   string
	stock int
}

func NewProduct(id, sku string, initialStock int) (*Product, error) {
	if id == "" || sku == "" {
		return nil, ErrInvalidSKU
	}
	if initialStock < 0 {
		return nil, ErrInvalidQuantity
	}
	return &Product{
		id:    id,
		sku:   sku,
		stock: initialStock,
	}, nil
}

func (p *Product) ID() string  { return p.id }
func (p *Product) SKU() string { return p.sku }
func (p *Product) Stock() int  { return p.stock }

func (p *Product) DeductStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	if p.stock < quantity {
		return ErrInsufficientStock
	}
	p.stock -= quantity
	return nil
}

func (p *Product) AddStock(quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}
	p.stock += quantity
	return nil
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
	mu       sync.RWMutex
	products map[string]*Product
}

func NewInMemoryProductRepo() *InMemoryProductRepo {
	return &InMemoryProductRepo{
		products: make(map[string]*Product),
	}
}

func (r *InMemoryProductRepo) GetBySKU(_ context.Context, sku string) (*Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.products[sku]
	if !exists {
		return nil, ErrProductNotFound
	}
	// Return defensive copy
	clone := *p
	return &clone, nil
}

func (r *InMemoryProductRepo) Save(_ context.Context, product *Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Store defensive copy
	clone := *product
	r.products[product.SKU()] = &clone
	return nil
}

// MockAlertPublisher is an adapter implementing AlertPublisher
type MockAlertPublisher struct {
	mu     sync.Mutex
	alerts map[string]int
}

func NewMockAlertPublisher() *MockAlertPublisher {
	return &MockAlertPublisher{
		alerts: make(map[string]int),
	}
}

func (m *MockAlertPublisher) PublishLowStock(_ context.Context, sku string, remaining int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alerts[sku] = remaining
	return nil
}

func (m *MockAlertPublisher) RecordedAlerts() map[string]int {
	m.mu.Lock()
	defer m.mu.Unlock()
	copyMap := make(map[string]int, len(m.alerts))
	for k, v := range m.alerts {
		copyMap[k] = v
	}
	return copyMap
}
