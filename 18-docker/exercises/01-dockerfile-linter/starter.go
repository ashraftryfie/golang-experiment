package dockerlinter

import (
	"errors"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

type LintViolation struct {
	RuleID  string
	Message string
}

func LintDockerfile(content string) ([]LintViolation, error) {
	return nil, ErrNotImplemented
}
