package store

import (
	"errors"
	"sync"

	"github.com/ashraftryfie/golang-experiment/06-packages-modules/challenges/modular-monolith/internal/domain"
)

var ErrNotFound = errors.New("product not found")

type Store interface {
	Save(p domain.Product) error
	Get(id string) (domain.Product, error)
	List() []domain.Product
}

type MemoryStore struct {
	mu       sync.RWMutex
	products map[string]domain.Product
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		products: make(map[string]domain.Product),
	}
}

func (s *MemoryStore) Save(p domain.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.products[p.ID] = p
	return nil
}

func (s *MemoryStore) Get(id string) (domain.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, exists := s.products[id]
	if !exists {
		return domain.Product{}, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) List() []domain.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]domain.Product, 0, len(s.products))
	for _, p := range s.products {
		res = append(res, p)
	}
	return res
}
