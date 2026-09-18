package crud_solution

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBookCRUDOperationsSolution(t *testing.T) {
	store := NewBookStore()
	mux, err := NewBookRouter(store)
	if err != nil {
		t.Fatalf("unexpected error creating router: %v", err)
	}

	// 1. Initial list empty array
	reqList := httptest.NewRequest(http.MethodGet, "/books", nil)
	recList := httptest.NewRecorder()
	mux.ServeHTTP(recList, reqList)

	if recList.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recList.Code)
	}

	// 2. Create Book (POST /books) -> 201 Created
	payload := `{"id":"book-101","title":"The Go Programming Language","author":"Alan A. A. Donovan"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(payload))
	recCreate := httptest.NewRecorder()
	mux.ServeHTTP(recCreate, reqCreate)

	if recCreate.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recCreate.Code)
	}
	if recCreate.Header().Get("Location") != "/books/book-101" {
		t.Errorf("expected Location /books/book-101, got %q", recCreate.Header().Get("Location"))
	}

	// 3. Fetch Book (GET /books/{id}) -> 200 OK
	reqGet := httptest.NewRequest(http.MethodGet, "/books/book-101", nil)
	recGet := httptest.NewRecorder()
	mux.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", recGet.Code)
	}
	var fetched Book
	if err := json.Unmarshal(recGet.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("failed to decode book: %v", err)
	}
	if fetched.Title != "The Go Programming Language" {
		t.Errorf("title mismatch: %q", fetched.Title)
	}

	// 4. Delete Book (DELETE /books/{id}) -> 204 No Content
	reqDel := httptest.NewRequest(http.MethodDelete, "/books/book-101", nil)
	recDel := httptest.NewRecorder()
	mux.ServeHTTP(recDel, reqDel)

	if recDel.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", recDel.Code)
	}

	// 5. Subsequent GET should return 404 Not Found
	reqGetMissing := httptest.NewRequest(http.MethodGet, "/books/book-101", nil)
	recGetMissing := httptest.NewRecorder()
	mux.ServeHTTP(recGetMissing, reqGetMissing)

	if recGetMissing.Code != http.StatusNotFound {
		t.Errorf("expected status 404 after deletion, got %d", recGetMissing.Code)
	}
}
