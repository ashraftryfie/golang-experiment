package crud_solution

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type BookStore struct {
	mu    sync.RWMutex
	books map[string]Book
}

func NewBookStore() *BookStore {
	return &BookStore{
		books: make(map[string]Book),
	}
}

func (s *BookStore) List() []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		res = append(res, b)
	}
	return res
}

func (s *BookStore) Get(id string) (Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	return b, ok
}

func (s *BookStore) Save(b Book) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.books[b.ID] = b
}

func (s *BookStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.books[id]; !ok {
		return false
	}
	delete(s.books, id)
	return true
}

func NewBookRouter(store *BookStore) (*http.ServeMux, error) {
	mux := http.NewServeMux()

	// GET /books
	mux.HandleFunc("GET /books", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(store.List())
	})

	// GET /books/{id}
	mux.HandleFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		book, found := store.Get(id)
		if !found {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(book)
	})

	// POST /books
	mux.HandleFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
		var book Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		if book.ID == "" {
			book.ID = fmt.Sprintf("book-%d", len(store.List())+1)
		}
		store.Save(book)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/books/"+book.ID)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(book)
	})

	// DELETE /books/{id}
	mux.HandleFunc("DELETE /books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if !store.Delete(id) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "not found"})
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return mux, nil
}
