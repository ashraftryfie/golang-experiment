# Exercise 03: Lightweight Forwarding Reverse Proxy (Tier 3 - Hard)

## 🎯 Problem Statement
In microservice architectures and API gateways, reverse proxies receive client requests, forward them to internal target services, and relay responses back to clients.

Implement `NewProxyHandler`:
```go
func NewProxyHandler(targetBaseURL string) http.Handler
```

### Requirements
1. Construct a new outbound HTTP request with the client's method (`r.Method`), target URL (`targetBaseURL + r.URL.Path`), and body (`r.Body`).
2. Forward headers from the client request to the outbound request (except hop-by-hop headers if any).
3. Execute the outbound request using `http.DefaultClient` or a configured `http.Client`.
4. Ensure outbound `resp.Body` is closed (`defer resp.Body.Close()`).
5. Copy response headers from upstream to client response writer `w`.
6. Write upstream HTTP status code `w.WriteHeader(resp.StatusCode)`.
7. Stream response body from `resp.Body` to `w` using `io.Copy`.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `NewProxyHandler`.
3. Run tests:
   ```powershell
   go test -v ./09-http/exercises/03-reverse-proxy
   ```
