# Exercise 03: RFC 7807 Problem Details Error Formatting (Tier 3 - Hard)

## 🎯 Problem Statement
Rather than returning arbitrary or inconsistent error payloads, modern microservices standardize error responses using **RFC 7807: Problem Details for HTTP APIs** with media type `application/problem+json`.

Implement the RFC 7807 writer:
```go
type InvalidParam struct {
    Name   string `json:"name"`
    Reason string `json:"reason"`
}

type ProblemDetails struct {
    Type          string         `json:"type"`
    Title         string         `json:"title"`
    Status        int            `json:"status"`
    Detail        string         `json:"detail,omitempty"`
    Instance      string         `json:"instance,omitempty"`
    InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
}

func WriteProblem(w http.ResponseWriter, p ProblemDetails) error
```

### Requirements
1. Set the response header `Content-Type: application/problem+json`.
2. If `p.Type` is empty, default it to `"about:blank"`.
3. Call `w.WriteHeader(p.Status)`.
4. Serialize `p` as JSON to `w` using `json.NewEncoder(w).Encode(p)`.
5. Return any error encountered during JSON encoding.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `WriteProblem`.
3. Run tests:
   ```powershell
   go test -v ./10-rest-api/exercises/03-rfc7807-errors
   ```
