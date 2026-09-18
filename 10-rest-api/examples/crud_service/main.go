package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Item struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

// MemoryStore provides thread-safe in-memory CRUD operations.
type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]Item),
	}
}

func (s *MemoryStore) Get(id string) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	return item, ok
}

func (s *MemoryStore) Set(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item.ID] = item
}

func (s *MemoryStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false
	}
	delete(s.items, id)
	return true
}

func (s *MemoryStore) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Item, 0, len(s.items))
	for _, it := range s.items {
		list = append(list, it)
	}
	return list
}

func main() {
	store := NewMemoryStore()
	mux := http.NewServeMux()

	// GET /items -> 200 OK
	mux.HandleFunc("GET /items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(store.List())
	})

	// GET /items/{id} -> 200 OK or 404 Not Found
	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		item, found := store.Get(id)
		if !found {
			http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(item)
	})

	// POST /items -> 201 Created
	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		var item Item
		if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if item.ID == "" {
			item.ID = fmt.Sprintf("item-%d", len(store.List())+1)
		}
		store.Set(item)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/items/"+item.ID)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(item)
	})

	// DELETE /items/{id} -> 204 No Content
	mux.HandleFunc("DELETE /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if !store.Delete(id) {
			http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	fmt.Println("REST CRUD Service example ready.")
}
