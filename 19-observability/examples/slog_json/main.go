package main

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type Account struct {
	ID       string
	Email    string
	Password string
	Token    string
}

// LogValue implements slog.LogValuer to automatically mask credentials
func (a Account) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", a.ID),
		slog.String("email", a.Email),
		slog.String("password", "[REDACTED]"),
		slog.String("token", "[REDACTED]"),
	)
}

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	acc := Account{
		ID:       "acc-987",
		Email:    "ops@example.com",
		Password: "plain-secret-password-123",
		Token:    "jwt-bearer-xyz-secret",
	}

	ctx := context.Background()

	logger.InfoContext(ctx, "account successfully authenticated",
		slog.Any("account", acc),
		slog.Duration("latency", 12*time.Millisecond),
		slog.String("client_ip", "192.168.1.10"),
	)
}
