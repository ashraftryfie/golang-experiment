package dockerlinter

import (
	"testing"
)

func TestDockerfileLinter_Valid(t *testing.T) {
	dockerfile := `
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /app/server /app/server
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
`
	violations, err := LintDockerfile(dockerfile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(violations) != 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func TestDockerfileLinter_Violations(t *testing.T) {
	// Anti-pattern Dockerfile: single stage, no CGO disabled, root user, bad copy order
	dockerfile := `
FROM golang:1.24
WORKDIR /app
COPY . .
COPY go.mod ./
RUN go build -o server .
CMD ["/app/server"]
`
	violations, err := LintDockerfile(dockerfile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedRules := map[string]bool{
		RuleMultiStage:      false,
		RuleCGODisabled:     false,
		RuleDependencyCache: false,
		RuleNonRootUser:     false,
	}

	for _, v := range violations {
		if _, ok := expectedRules[v.RuleID]; ok {
			expectedRules[v.RuleID] = true
		}
	}

	for rule, found := range expectedRules {
		if !found {
			t.Errorf("expected violation for rule %s, but none was reported", rule)
		}
	}
}
