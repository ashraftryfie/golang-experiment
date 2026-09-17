package main

import (
	"testing"
)

func TestMiddlewarePipeline(t *testing.T) {
	logs := make([]string, 0)
	logger := func(msg string) {
		logs = append(logs, msg)
	}

	handler := func(req Request) Response {
		return Response{StatusCode: 200, Body: "Success: " + req.Body}
	}

	pipeline := Chain(
		LoggingMiddleware(logger),
		AuthMiddleware("bearer-token"),
	)(handler)

	// Test 1: Successful Auth & Execution
	res := pipeline(Request{Path: "/data", AuthToken: "bearer-token", Body: "payload"})
	if res.StatusCode != 200 || res.Body != "Success: payload" {
		t.Errorf("expected success, got %v", res)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 log entries, got %d", len(logs))
	}

	// Test 2: Auth Rejection
	logs = logs[:0]
	badRes := pipeline(Request{Path: "/data", AuthToken: "invalid", Body: "payload"})
	if badRes.StatusCode != 401 {
		t.Errorf("expected 401, got %d", badRes.StatusCode)
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	panickingHandler := func(req Request) Response {
		panic("database connection dropped!")
	}

	recoveredPipeline := Chain(
		RecoveryMiddleware(nil),
	)(panickingHandler)

	res := recoveredPipeline(Request{Path: "/crash"})
	if res.StatusCode != 500 || res.Body != "Internal Server Error" {
		t.Errorf("expected recovered 500 error, got %v", res)
	}
}
