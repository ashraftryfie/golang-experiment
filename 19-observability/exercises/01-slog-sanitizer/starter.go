package slogsanitizer

import (
	"context"
	"log/slog"
)

type SanitizingHandler struct {
	next slog.Handler
}

func NewSanitizingHandler(underlying slog.Handler) *SanitizingHandler {
	return nil
}

func (h *SanitizingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return false
}

func (h *SanitizingHandler) Handle(ctx context.Context, r slog.Record) error {
	return nil
}

func (h *SanitizingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *SanitizingHandler) WithGroup(name string) slog.Handler {
	return h
}
