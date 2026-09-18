package composevalidator

import (
	"errors"
)

var ErrNotImplemented = errors.New("exercise not implemented yet")

type ComposeViolation struct {
	RuleID  string
	Service string
	Message string
}

func ValidateCompose(content string, requiredServices []string) ([]ComposeViolation, error) {
	return nil, ErrNotImplemented
}
