# Exercise 03: W3C Distributed Trace Propagator

## Objective
Implement a parser and propagator for the W3C Trace Context standard (`traceparent` header) to enable distributed request correlation across microservice boundaries.

## Requirements
1. **Traceparent Structure**:
   `00-<32-hex-trace-id>-<16-hex-parent-id>-<2-hex-flags>`
2. **Validation Rules**:
   - Version must be `00`.
   - Trace ID must be 32 lowercase hex characters and non-zero (`ErrInvalidTraceID`).
   - Parent ID must be 16 lowercase hex characters and non-zero (`ErrInvalidParentID`).
   - Flags must be 2 lowercase hex characters.
3. **HTTP Helpers**:
   - `InjectTraceparent(req *http.Request, tc TraceContext)` sets `traceparent` header.
   - `ExtractTraceparent(r *http.Request) (*TraceContext, bool)` extracts and parses header.
   - `GenerateNewTraceContext() TraceContext` generates random cryptographic IDs with flag `01` (recorded).

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./19-observability/exercises/03-w3c-trace-propagator`
