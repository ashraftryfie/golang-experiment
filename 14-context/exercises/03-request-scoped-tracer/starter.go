package tracer

import (
	"context"
	"errors"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

type TraceInfo struct {
	TraceID   string
	RequestID string
	UserID    string
	Tags      map[string]string
}

// WithTraceInfo attaches TraceInfo to the context using a private key.
func WithTraceInfo(ctx context.Context, info TraceInfo) context.Context {
	return nil
}

// GetTraceInfo extracts TraceInfo from the context. Returns false if absent.
func GetTraceInfo(ctx context.Context) (TraceInfo, bool) {
	return TraceInfo{}, false
}

// AddTraceTag creates a new child context with the added tag without mutating the parent.
func AddTraceTag(ctx context.Context, key, value string) (context.Context, error) {
	return nil, ErrNotImplemented
}
