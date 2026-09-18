package dockerlinter

import (
	"errors"
	"testing"
)

func TestDockerfileLinter_Starter(t *testing.T) {
	dockerfile := `
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/server /app/server
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
`
	violations, err := LintDockerfile(dockerfile)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: LintDockerfile")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Errorf("expected 0 violations for valid dockerfile, got %d", len(violations))
	}
}
