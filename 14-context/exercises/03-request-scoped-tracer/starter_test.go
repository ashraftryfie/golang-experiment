package tracer

import (
	"context"
	"errors"
	"testing"
)

func TestTracer_Starter(t *testing.T) {
	info := TraceInfo{
		TraceID:   "trc-123",
		RequestID: "req-456",
		UserID:    "usr-789",
		Tags:      map[string]string{"env": "prod"},
	}

	ctx := context.Background()
	ctxWithTrace := WithTraceInfo(ctx, info)
	if ctxWithTrace == nil {
		t.Skip("Skipping unimplemented exercise: WithTraceInfo")
	}

	extracted, ok := GetTraceInfo(ctxWithTrace)
	if !ok {
		t.Fatalf("expected to extract trace info")
	}
	if extracted.TraceID != "trc-123" {
		t.Errorf("expected trace ID trc-123, got %s", extracted.TraceID)
	}

	_, err := AddTraceTag(ctxWithTrace, "version", "v1.2.0")
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: AddTraceTag")
	}
}
