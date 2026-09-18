# Challenge: Production-Ready API Gateway Router

## 🎯 Objective
Construct a lightweight API Gateway router combining Go 1.22+ method routing, path parameter extraction, middleware chains (Request ID & API Key verification), and unit tests using `net/http/httptest`.

## 📋 Requirements
1. **Public Routes**:
   - `GET /healthz`: Returns status `200 OK` and JSON `{"status":"ok"}`. Open to all clients without authentication.
2. **Protected Routes**:
   - `GET /v1/services/{name}`: Requires valid `X-API-Key`.
   - Returns status `200 OK` and JSON `{"service":"<name>","status":"active"}`.
3. **Gateway Middlewares**:
   - **Request ID**: Ensures `X-Request-ID` is present on both request context and response header.
   - **Auth**: If `r.URL.Path` begins with `/v1/`, requires `X-API-Key == gatewayKey`. Missing or invalid key yields `401 Unauthorized` with JSON `{"error":"unauthorized"}`.
4. **Zero Third-Party Dependencies**: Pure standard library (`net/http`, `net/http/httptest`, `encoding/json`).

## 🧪 Verification
```powershell
go test -v ./09-http/challenges/api-gateway/...
```
