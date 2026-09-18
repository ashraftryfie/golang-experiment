package slogsanitizer

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestSanitizer_RedactsSensitiveFields(t *testing.T) {
	buf := new(bytes.Buffer)
	base := slog.NewJSONHandler(buf, nil)
	sanitized := NewSanitizingHandler(base)
	logger := slog.New(sanitized)

	logger.InfoContext(context.Background(), "auth event",
		slog.String("username", "john_doe"),
		slog.String("PASSWORD", "mysecretpwd!"),
		slog.String("API_KEY", "key-live-12345"),
		slog.Group("payment",
			slog.String("provider", "stripe"),
			slog.String("credit_card", "4111-2222-3333-4444"),
		),
	)

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed parsing json log: %v", err)
	}

	if logEntry["username"] != "john_doe" {
		t.Errorf("expected username john_doe, got %v", logEntry["username"])
	}

	if logEntry["PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected PASSWORD to be [REDACTED], got %v", logEntry["PASSWORD"])
	}

	if logEntry["API_KEY"] != "[REDACTED]" {
		t.Errorf("expected API_KEY to be [REDACTED], got %v", logEntry["API_KEY"])
	}

	paymentGroup, ok := logEntry["payment"].(map[string]any)
	if !ok {
		t.Fatalf("expected payment group in log entry")
	}

	if paymentGroup["provider"] != "stripe" {
		t.Errorf("expected payment.provider to be stripe, got %v", paymentGroup["provider"])
	}

	if paymentGroup["credit_card"] != "[REDACTED]" {
		t.Errorf("expected payment.credit_card to be [REDACTED], got %v", paymentGroup["credit_card"])
	}
}
