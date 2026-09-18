# Exercise 01: Slog Sensitive Attribute Sanitizer

## Objective
Implement a custom `slog.Handler` that transparently intercepts structured log records and redacts sensitive PII (passwords, tokens, API keys, secrets) before writing them to downstream sinks.

## Requirements
1. **`NewSanitizingHandler(underlying slog.Handler) slog.Handler`**:
   - Wraps any standard `slog.Handler` (e.g. `slog.JSONHandler` or `slog.TextHandler`).
2. **Key Matching**:
   - Any attribute whose key matches (case-insensitive) `password`, `token`, `secret`, `authorization`, `api_key`, or `credit_card` must have its value replaced with `"[REDACTED]"`.
3. **Recursive Group Sanitization**:
   - Nested `slog.Group` structures must be traversed and sanitized recursively.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./19-observability/exercises/01-slog-sanitizer`
