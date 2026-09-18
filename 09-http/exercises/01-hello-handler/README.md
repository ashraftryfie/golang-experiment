# Exercise 01: Go 1.22+ Greeting Handler (Tier 1 - Easy)

## 🎯 Problem Statement
Implement a standard Go HTTP router using Go 1.22+ path wildcards:
```go
func NewGreetingRouter() *http.ServeMux
```

### Requirements
1. Register `GET /greet/{name}` on the returned `*http.ServeMux`.
2. Extract the name parameter using `r.PathValue("name")`.
3. If `name` is empty, return status code `400 Bad Request` with plain text or JSON error.
4. If `name` is present, set header `Content-Type: application/json`, status `200 OK`, and return:
   ```json
   {
     "message": "Hello, <name>!",
     "status": "success"
   }
   ```

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `NewGreetingRouter`.
3. Run tests:
   ```powershell
   go test -v ./09-http/exercises/01-hello-handler
   ```
