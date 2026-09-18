package logscan

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var ErrEmptyTargetLevel = errors.New("target level cannot be empty")

func FilterLogs(r io.Reader, targetLevel string) ([]string, error) {
	if targetLevel == "" {
		return nil, ErrEmptyTargetLevel
	}

	var results []string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, targetLevel) {
			results = append(results, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
