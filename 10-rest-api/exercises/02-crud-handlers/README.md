# Exercise 02: Thread-Safe REST CRUD Handlers (Tier 2 - Medium)

## 🎯 Problem Statement
Implement a thread-safe RESTful CRUD router for a book inventory service using Go standard library `sync.RWMutex` and Go 1.22+ `http.ServeMux`.

```go
type Book struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

func NewBookRouter(store *BookStore) *http.ServeMux
```

### Endpoints & Semantics
1. `GET /books`:
   - Returns status `200 OK`.
   - Content-Type `application/json`.
   - Body: JSON array of all books (return `[]` when empty, not `null`).
2. `GET /books/{id}`:
   - If found: status `200 OK`, JSON book payload.
   - If not found: status `404 Not Found`, JSON `{"error":"not found"}`.
3. `POST /books`:
   - Decodes JSON body.
   - If `ID` is empty or missing, generate a unique ID.
   - Stores book in `BookStore`.
   - Header `Location: /books/{id}`.
   - Status `201 Created`, JSON book payload.
4. `DELETE /books/{id}`:
   - If found: removes book, returns status `204 No Content` (empty body).
   - If not found: returns status `404 Not Found`.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `BookStore` methods and `NewBookRouter`.
3. Run tests:
   ```powershell
   go test -v ./10-rest-api/exercises/02-crud-handlers
   ```
