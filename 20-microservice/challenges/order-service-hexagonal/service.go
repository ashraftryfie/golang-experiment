package orderservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ==========================================
// 1. DOMAIN LAYER
// ==========================================

var (
	ErrEmptyCustomerID       = errors.New("customer_id cannot be empty")
	ErrEmptyOrderItems       = errors.New("order must contain at least one item")
	ErrInvalidItemQuantity   = errors.New("item quantity must be greater than zero")
	ErrInvalidItemPrice      = errors.New("item price cannot be negative")
	ErrOrderNotFound         = errors.New("order not found")
	ErrAlreadyPaid           = errors.New("order is already paid")
	ErrAlreadyCancelled      = errors.New("order is already cancelled")
	ErrCannotCancelShipped   = errors.New("cannot cancel an order that has already shipped")
	ErrPaymentFailed         = errors.New("payment processing failed")
	ErrInvalidTransition    = errors.New("invalid order status transition")
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusPaid      OrderStatus = "PAID"
	StatusShipped   OrderStatus = "SHIPPED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type OrderItem struct {
	ProductID      string `json:"product_id"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
}

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
}

func NewOrder(id, customerID string, items []OrderItem) (*Order, error) {
	if strings.TrimSpace(customerID) == "" {
		return nil, ErrEmptyCustomerID
	}
	if len(items) == 0 {
		return nil, ErrEmptyOrderItems
	}
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, ErrInvalidItemQuantity
		}
		if item.UnitPriceCents < 0 {
			return nil, ErrInvalidItemPrice
		}
	}

	itemsCopy := make([]OrderItem, len(items))
	copy(itemsCopy, items)

	return &Order{
		ID:         id,
		CustomerID: customerID,
		Items:      itemsCopy,
		Status:     StatusPending,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

func (o *Order) TotalCents() int64 {
	var total int64
	for _, item := range o.Items {
		total += int64(item.Quantity) * item.UnitPriceCents
	}
	return total
}

func (o *Order) MarkPaid() error {
	if o.Status == StatusPaid {
		return ErrAlreadyPaid
	}
	if o.Status == StatusCancelled {
		return ErrAlreadyCancelled
	}
	if o.Status != StatusPending {
		return fmt.Errorf("%w: cannot mark as paid from status %s", ErrInvalidTransition, o.Status)
	}
	o.Status = StatusPaid
	return nil
}

func (o *Order) MarkShipped() error {
	if o.Status != StatusPaid {
		return fmt.Errorf("%w: order must be paid before shipping, current status: %s", ErrInvalidTransition, o.Status)
	}
	o.Status = StatusShipped
	return nil
}

func (o *Order) Cancel() error {
	if o.Status == StatusShipped {
		return ErrCannotCancelShipped
	}
	if o.Status == StatusCancelled {
		return ErrAlreadyCancelled
	}
	o.Status = StatusCancelled
	return nil
}

// ==========================================
// 2. PORTS
// ==========================================

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id string) (*Order, error)
	ListByCustomer(ctx context.Context, customerID string) ([]*Order, error)
}

type PaymentGateway interface {
	Charge(ctx context.Context, customerID string, amountCents int64) (string, error)
}

type EventDispatcher interface {
	Dispatch(ctx context.Context, eventName string, payload any) error
}

// ==========================================
// 3. SERVICE LAYER
// ==========================================

type OrderService struct {
	repo       OrderRepository
	payment    PaymentGateway
	dispatcher EventDispatcher
	idGen      func() string
}

func NewOrderService(repo OrderRepository, payment PaymentGateway, dispatcher EventDispatcher) *OrderService {
	return &OrderService{
		repo:       repo,
		payment:    payment,
		dispatcher: dispatcher,
		idGen: func() string {
			return fmt.Sprintf("ord-%d", time.Now().UnixNano())
		},
	}
}

// SetIDGenerator allows injecting a custom ID generator for deterministic testing
func (s *OrderService) SetIDGenerator(fn func() string) {
	s.idGen = fn
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, items []OrderItem) (*Order, error) {
	orderID := s.idGen()
	order, err := NewOrder(orderID, customerID, items)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("saving order: %w", err)
	}

	if s.dispatcher != nil {
		_ = s.dispatcher.Dispatch(ctx, "order.created", order)
	}

	return order, nil
}

func (s *OrderService) PayOrder(ctx context.Context, orderID string) error {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := order.MarkPaid(); err != nil {
		return err
	}

	txID, err := s.payment.Charge(ctx, order.CustomerID, order.TotalCents())
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPaymentFailed, err)
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("saving paid order: %w", err)
	}

	if s.dispatcher != nil {
		_ = s.dispatcher.Dispatch(ctx, "order.paid", map[string]any{
			"order_id":       order.ID,
			"transaction_id": txID,
			"amount_cents":   order.TotalCents(),
		})
	}

	return nil
}

func (s *OrderService) ShipOrder(ctx context.Context, orderID string) error {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := order.MarkShipped(); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("saving shipped order: %w", err)
	}

	if s.dispatcher != nil {
		_ = s.dispatcher.Dispatch(ctx, "order.shipped", order)
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID string) error {
	order, err := s.repo.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if err := order.Cancel(); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, order); err != nil {
		return fmt.Errorf("saving cancelled order: %w", err)
	}

	if s.dispatcher != nil {
		_ = s.dispatcher.Dispatch(ctx, "order.cancelled", order)
	}

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	return s.repo.GetByID(ctx, orderID)
}

// ==========================================
// 4. ADAPTERS
// ==========================================

// InMemoryOrderRepository is a concurrency-safe repository adapter
type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders map[string]*Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{orders: make(map[string]*Order)}
}

func (r *InMemoryOrderRepository) Save(_ context.Context, order *Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	itemsCopy := make([]OrderItem, len(order.Items))
	copy(itemsCopy, order.Items)

	clone := *order
	clone.Items = itemsCopy
	r.orders[order.ID] = &clone
	return nil
}

func (r *InMemoryOrderRepository) GetByID(_ context.Context, id string) (*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, ErrOrderNotFound
	}

	itemsCopy := make([]OrderItem, len(order.Items))
	copy(itemsCopy, order.Items)

	clone := *order
	clone.Items = itemsCopy
	return &clone, nil
}

func (r *InMemoryOrderRepository) ListByCustomer(_ context.Context, customerID string) ([]*Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*Order
	for _, order := range r.orders {
		if order.CustomerID == customerID {
			itemsCopy := make([]OrderItem, len(order.Items))
			copy(itemsCopy, order.Items)
			clone := *order
			clone.Items = itemsCopy
			result = append(result, &clone)
		}
	}
	return result, nil
}

// MockPaymentGateway is an adapter implementing PaymentGateway
type MockPaymentGateway struct {
	mu          sync.Mutex
	ShouldFail  bool
	Charges     []int64
	NextTxIDGen func() string
}

func NewMockPaymentGateway() *MockPaymentGateway {
	return &MockPaymentGateway{
		NextTxIDGen: func() string {
			return fmt.Sprintf("tx-%d", time.Now().UnixNano())
		},
	}
}

func (m *MockPaymentGateway) Charge(_ context.Context, _ string, amountCents int64) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ShouldFail {
		return "", errors.New("card declined")
	}
	m.Charges = append(m.Charges, amountCents)
	return m.NextTxIDGen(), nil
}

// MockEventDispatcher is an adapter recording dispatched domain events
type MockEventDispatcher struct {
	mu     sync.Mutex
	Events []string
}

func NewMockEventDispatcher() *MockEventDispatcher {
	return &MockEventDispatcher{}
}

func (m *MockEventDispatcher) Dispatch(_ context.Context, eventName string, _ any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Events = append(m.Events, eventName)
	return nil
}

// OrderHTTPHandler is a primary adapter translating HTTP requests to Service calls
type OrderHTTPHandler struct {
	service *OrderService
	mux     *http.ServeMux
}

func NewOrderHTTPHandler(service *OrderService) *OrderHTTPHandler {
	h := &OrderHTTPHandler{
		service: service,
		mux:     http.NewServeMux(),
	}
	h.routes()
	return h
}

func (h *OrderHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *OrderHTTPHandler) routes() {
	h.mux.HandleFunc("POST /orders", h.handleCreateOrder)
	h.mux.HandleFunc("POST /orders/{id}/pay", h.handlePayOrder)
	h.mux.HandleFunc("GET /orders/{id}", h.handleGetOrder)
}

type createOrderRequest struct {
	CustomerID string      `json:"customer_id"`
	Items      []OrderItem `json:"items"`
}

func (h *OrderHTTPHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateOrder(r.Context(), req.CustomerID, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(order)
}

func (h *OrderHTTPHandler) handlePayOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		http.Error(w, "missing order id", http.StatusBadRequest)
		return
	}

	if err := h.service.PayOrder(r.Context(), orderID); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"PAID"}`))
}

func (h *OrderHTTPHandler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("id")
	if orderID == "" {
		http.Error(w, "missing order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}
