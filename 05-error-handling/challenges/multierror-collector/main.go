package main

import (
	"errors"
	"fmt"
	"strings"
)

type MultiError struct {
	errors []error
}

func (m *MultiError) Add(err error) {
	if err != nil {
		m.errors = append(m.errors, err)
	}
}

func (m *MultiError) ErrorOrNil() error {
	if m == nil || len(m.errors) == 0 {
		return nil
	}
	return m
}

func (m *MultiError) Error() string {
	if len(m.errors) == 0 {
		return "no errors"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("encountered %d errors:\n", len(m.errors)))
	for _, e := range m.errors {
		sb.WriteString(fmt.Sprintf("  - %s\n", e.Error()))
	}
	return strings.TrimRight(sb.String(), "\n")
}

// Unwrap implements Go 1.20+ multi-error unwrapping
func (m *MultiError) Unwrap() []error {
	return m.errors
}

func (m *MultiError) Errors() []error {
	return m.errors
}

func main() {
	var errs MultiError
	errs.Add(errors.New("invalid email"))
	errs.Add(errors.New("password too weak"))

	if err := errs.ErrorOrNil(); err != nil {
		fmt.Println(err)
	}
}
