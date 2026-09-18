package slogsanitizer

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestSanitizer_Starter(t *testing.T) {
	buf := new(bytes.Buffer)
	baseHandler := slog.NewJSONHandler(buf, nil)
	sanitizedHandler := NewSanitizingHandler(baseHandler)
	if sanitizedHandler == nil {
		t.Skip("Skipping unimplemented exercise: NewSanitizingHandler")
	}

	logger := slog.New(sanitizedHandler)
	logger.InfoContext(context.Background(), "user login", slog.String("password", "supersecret123"))

	if strings.Contains(buf.String(), "supersecret123") {
		t.Fatalf("sensitive password was not redacted from log output: %s", buf.String())
	}
}
