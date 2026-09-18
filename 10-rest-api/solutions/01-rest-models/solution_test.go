package restmodels_solution

import (
	"testing"
)

func TestValidateSolution(t *testing.T) {
	tests := []struct {
		name       string
		req        CreateUserRequest
		wantField  string
		wantErrors int
	}{
		{
			name:       "valid request has 0 errors",
			req:        CreateUserRequest{Username: "gopher", Email: "gopher@golang.org", Age: 30},
			wantErrors: 0,
		},
		{
			name:       "username too short",
			req:        CreateUserRequest{Username: "go", Email: "gopher@golang.org", Age: 30},
			wantField:  "username",
			wantErrors: 1,
		},
		{
			name:       "username with space",
			req:        CreateUserRequest{Username: "go pher", Email: "gopher@golang.org", Age: 30},
			wantField:  "username",
			wantErrors: 1,
		},
		{
			name:       "invalid email",
			req:        CreateUserRequest{Username: "gopher", Email: "invalid-email", Age: 30},
			wantField:  "email",
			wantErrors: 1,
		},
		{
			name:       "age under 18",
			req:        CreateUserRequest{Username: "gopher", Email: "gopher@golang.org", Age: 16},
			wantField:  "age",
			wantErrors: 1,
		},
		{
			name:       "multiple invalid fields",
			req:        CreateUserRequest{Username: "a", Email: "bad", Age: 10},
			wantErrors: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.req.Validate()
			if len(result) != tc.wantErrors {
				t.Fatalf("got %d errors, want %d: %+v", len(result), tc.wantErrors, result)
			}
			if tc.wantField != "" && len(result) > 0 {
				found := false
				for _, err := range result {
					if err.Field == tc.wantField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected error on field %q, got %+v", tc.wantField, result)
				}
			}
		})
	}
}
