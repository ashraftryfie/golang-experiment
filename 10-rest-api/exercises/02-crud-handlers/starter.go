package crud

import (
	"errors"
	"net/http"
	"sync"
)

var ErrNotImplemented = errors.New("TODO: implement BookStore and NewBookRouter")

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

// NewBookRouter wires up REST routes on an http.ServeMux.
func NewBookRouter(store *BookStore) (*http.ServeMux, error) {
	// TODO: Register GET /books, GET /books/{id}, POST /books, DELETE /books/{id}
	return nil, ErrNotImplemented
}
