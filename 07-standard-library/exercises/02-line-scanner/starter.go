package logscan

import (
	"errors"
	"io"
)

var (
	ErrEmptyTargetLevel = errors.New("target level cannot be empty")
	ErrNotImplemented   = errors.New("TODO: implement")
)

// FilterLogs reads r line-by-line and returns all lines containing targetLevel.
func FilterLogs(r io.Reader, targetLevel string) ([]string, error) {
	// TODO: Validate targetLevel, initialize bufio.Scanner, collect matching lines
	return nil, ErrNotImplemented
}
