package main

import (
	"errors"
	"fmt"
)

// 1. Sentinel Errors (exported package variables)
var (
	ErrNotFound   = errors.New("resource not found")
	ErrPermission = errors.New("permission denied")
)

// 2. Custom Error Struct
type QueryError struct {
	Query string
	Err   error
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("failed to execute query %q: %v", e.Query, e.Err)
}

// Unwrap enables errors.Is and errors.As to inspect the inner error
func (e *QueryError) Unwrap() error {
	return e.Err
}

func findUser(id int) (*string, error) {
	if id <= 0 {
		return nil, &QueryError{
			Query: fmt.Sprintf("SELECT * FROM users WHERE id=%d", id),
			Err:   ErrNotFound,
		}
	}
	name := "Ashraf"
	return &name, nil
}

func main() {
	_, err := findUser(-1)
	if err != nil {
		fmt.Printf("Full error string: %v\n", err)

		// Inspect with errors.Is (unwraps QueryError to find ErrNotFound)
		if errors.Is(err, ErrNotFound) {
			fmt.Println("-> Detected ErrNotFound via errors.Is()")
		}

		// Inspect with errors.As to extract the query string
		var qErr *QueryError
		if errors.As(err, &qErr) {
			fmt.Printf("-> Extracted QueryError: Offending SQL was: %s\n", qErr.Query)
		}
	}
}
