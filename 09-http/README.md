# Stage 09: HTTP Fundamentals & Go 1.22+ Routing

Welcome to **Stage 09** of your Go engineering journey. In this stage, you master building HTTP servers, routing, middleware pipelines, and HTTP clients using solely the Go standard library (`net/http`, `net/http/httptest`).

---

## 🧭 Mental Model

The entire Go web ecosystem is built upon a single, elegant interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Any type that implements `ServeHTTP` can handle HTTP traffic. Functions with the signature `func(w http.ResponseWriter, r *http.Request)` can be converted into handlers via `http.HandlerFunc(myFunc)`.

### 1. Modern Go 1.22+ Enhanced `http.ServeMux`
Prior to Go 1.22, developers relied on third-party routers (`chi`, `gorilla/mux`) for HTTP method matching and path parameters. In Go 1.22+, the standard `http.ServeMux` supports both natively:

```go
mux := http.NewServeMux()

// Exact method + path variable
mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id") // Extract path parameter directly!
    w.Write([]byte("User ID: " + id))
})

// Subtree routing
mux.HandleFunc("GET /static/", staticHandler)
```

### 2. Middleware as Decorator Functions
A middleware in Go is a higher-order function that wraps an `http.Handler`:
```go
type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s finished in %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

### 3. Production Server Hygiene
Never use `http.ListenAndServe(":8080", nil)` in production! Default server configurations have no timeouts, making servers vulnerable to Slowloris attacks.
Always configure explicit timeouts on an `http.Server`:
```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

---

## 🎯 Learning Objectives

By the end of this stage, you will:
- [x] Implement custom HTTP handlers using both `http.Handler` and `http.HandlerFunc`.
- [x] Utilize Go 1.22+ `http.ServeMux` method matching and `r.PathValue(key)`.
- [x] Chain composable middlewares for logging, authentication, and error recovery.
- [x] Test HTTP handlers without opening real network sockets using `net/http/httptest`.
- [x] Configure production-grade timeouts on `http.Server`.

---

## 🗂️ Stage Structure

```
09-http/
├── README.md                           # Curriculum & mental model
├── examples/
│   ├── basic_server/main.go            # Go 1.22+ ServeMux with path params & timeouts
│   └── middleware_chain/main.go        # Composable middleware chaining
├── exercises/
│   ├── 01-hello-handler/               # JSON responses and PathValue routing
│   ├── 02-middleware-pipeline/         # API key auth & timing middleware
│   └── 03-reverse-proxy/               # Forwarding proxy with httptest verification
├── solutions/
│   ├── 01-hello-handler/               # Reference implementation
│   ├── 02-middleware-pipeline/         # Reference implementation
│   └── 03-reverse-proxy/               # Reference implementation
└── challenges/
    └── api-gateway/                    # Production router with auth, tracing & httptest
```

---

## ⚠️ Common Pitfalls & Gotchas

1. **Writing Status Code After Body**:
   ```go
   // ❌ WRONG: w.Write implicitly writes status 200 OK!
   w.Write([]byte("Not Found"))
   w.WriteHeader(http.StatusNotFound) // Warning: superfluous response.WriteHeader call

   // ✅ CORRECT: Always call WriteHeader before w.Write!
   w.WriteHeader(http.StatusNotFound)
   w.Write([]byte("Not Found"))
   ```

2. **Forgetting to `return` after `http.Error`**:
   ```go
   // ❌ WRONG: Code continues executing after writing error!
   if err != nil {
       http.Error(w, "Bad Request", http.StatusBadRequest)
   }
   doSensitiveOperation() // Runs even if request was bad!

   // ✅ CORRECT:
   if err != nil {
       http.Error(w, "Bad Request", http.StatusBadRequest)
       return
   }
   ```

3. **Leaking Client Response Bodies**:
   - Always close response bodies when using `http.Client`:
   ```go
   resp, err := client.Do(req)
   if err != nil { return err }
   defer resp.Body.Close() // MANDATORY to enable TCP connection reuse
   ```

---

## 🧪 Verification Commands

```powershell
# Run tests for this stage
go test -v ./09-http/...

# Format and vet
go vet ./09-http/...
go fmt ./09-http/...
```

---

## 🔗 Connections
- **Prerequisites**: [Stage 08: Files, JSON & Serialization](../08-files-json/README.md).
- **Next Stage**: [Stage 10: REST API Engineering](../10-rest-api/README.md) (CRUD architecture, request payload validation, REST standards).
