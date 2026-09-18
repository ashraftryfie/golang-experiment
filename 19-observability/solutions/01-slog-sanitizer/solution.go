package slogsanitizer

import (
	"context"
	"log/slog"
	"strings"
)

type SanitizingHandler struct {
	next slog.Handler
}

func NewSanitizingHandler(underlying slog.Handler) *SanitizingHandler {
	return &SanitizingHandler{next: underlying}
}

func (h *SanitizingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

var sensitiveKeys = map[string]bool{
	"password":      true,
	"token":         true,
	"secret":        true,
	"authorization": true,
	"api_key":       true,
	"credit_card":   true,
}

func sanitizeAttr(a slog.Attr) slog.Attr {
	// Handle groups recursively
	if a.Value.Kind() == slog.KindGroup {
		attrs := a.Value.Group()
		sanitized := make([]slog.Attr, len(attrs))
		for i, sub := range attrs {
			sanitized[i] = sanitizeAttr(sub)
		}
		return slog.Attr{
			Key:   a.Key,
			Value: slog.GroupValue(sanitized...),
		}
	}

	keyLower := strings.ToLower(a.Key)
	if sensitiveKeys[keyLower] {
		return slog.String(a.Key, "[REDACTED]")
	}

	return a
}

func (h *SanitizingHandler) Handle(ctx context.Context, r slog.Record) error {
	sanitizedRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)

	r.Attrs(func(a slog.Attr) bool {
		sanitizedRecord.AddAttrs(sanitizeAttr(a))
		return true
	})

	return h.next.Handle(ctx, sanitizedRecord)
}

func (h *SanitizingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	sanitized := make([]slog.Attr, len(attrs))
	for i, a := range attrs {
		sanitized[i] = sanitizeAttr(a)
	}
	return &SanitizingHandler{next: h.next.WithAttrs(sanitized)}
}

func (h *SanitizingHandler) WithGroup(name string) slog.Handler {
	return &SanitizingHandler{next: h.next.WithGroup(name)}
}
