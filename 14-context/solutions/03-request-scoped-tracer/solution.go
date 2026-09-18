package tracer

import (
	"context"
	"errors"
)

var ErrNoTraceInfo = errors.New("trace info not found in context")

// traceContextKey is an unexported type to prevent key collisions across packages.
type traceContextKey struct{}

var traceKey = traceContextKey{}

type TraceInfo struct {
	TraceID   string
	RequestID string
	UserID    string
	Tags      map[string]string
}

// cloneTags creates a deep copy of a tag map to prevent race conditions.
func cloneTags(tags map[string]string) map[string]string {
	if tags == nil {
		return make(map[string]string)
	}
	cloned := make(map[string]string, len(tags)+1)
	for k, v := range tags {
		cloned[k] = v
	}
	return cloned
}

// WithTraceInfo attaches TraceInfo to the context using a private key.
func WithTraceInfo(ctx context.Context, info TraceInfo) context.Context {
	// Deep copy tags to protect against caller mutating input map later
	safeInfo := TraceInfo{
		TraceID:   info.TraceID,
		RequestID: info.RequestID,
		UserID:    info.UserID,
		Tags:      cloneTags(info.Tags),
	}
	return context.WithValue(ctx, traceKey, safeInfo)
}

// GetTraceInfo extracts TraceInfo from the context. Returns false if absent.
func GetTraceInfo(ctx context.Context) (TraceInfo, bool) {
	val := ctx.Value(traceKey)
	if val == nil {
		return TraceInfo{}, false
	}

	info, ok := val.(TraceInfo)
	if !ok {
		return TraceInfo{}, false
	}

	// Return a copy with cloned tags so caller cannot mutate stored map
	return TraceInfo{
		TraceID:   info.TraceID,
		RequestID: info.RequestID,
		UserID:    info.UserID,
		Tags:      cloneTags(info.Tags),
	}, true
}

// AddTraceTag creates a new child context with the added tag without mutating the parent.
func AddTraceTag(ctx context.Context, key, value string) (context.Context, error) {
	info, ok := GetTraceInfo(ctx)
	if !ok {
		return nil, ErrNoTraceInfo
	}

	info.Tags[key] = value
	return WithTraceInfo(ctx, info), nil
}
