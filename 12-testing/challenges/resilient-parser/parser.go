package kvparser

import (
	"errors"
	"fmt"
	"strings"
)

var ErrMissingDelimiter = errors.New("missing '=' delimiter in pair")

// ParseKV safely parses configuration key-value pairs separated by '&' or ';'.
func ParseKV(input string) (map[string]string, error) {
	result := make(map[string]string)
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return result, nil
	}

	// Split by both '&' and ';'
	normalized := strings.ReplaceAll(trimmed, ";", "&")
	pairs := strings.Split(normalized, "&")

	for _, rawPair := range pairs {
		pair := strings.TrimSpace(rawPair)
		if pair == "" || strings.HasPrefix(pair, "#") {
			continue // skip empty chunks or comments
		}

		idx := strings.IndexByte(pair, '=')
		if idx == -1 {
			return nil, fmt.Errorf("%w: %q", ErrMissingDelimiter, pair)
		}

		key := strings.TrimSpace(pair[:idx])
		val := strings.TrimSpace(pair[idx+1:])

		if key == "" {
			return nil, errors.New("empty key in pair")
		}

		result[key] = val
	}

	return result, nil
}
