package tabletests

import (
	"errors"
	"strings"
	"testing"
)

func ValidatePassword(pw string) error {
	if len(pw) < 8 {
		return errors.New("password too short: minimum 8 characters")
	}
	if !strings.ContainsAny(pw, "0123456789") {
		return errors.New("password must contain at least one digit")
	}
	if !strings.ContainsAny(pw, "!@#$%^&*") {
		return errors.New("password must contain at least one special character")
	}
	return nil
}

func assertErrorSubstr(t *testing.T, err error, substr string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error containing %q, got nil", substr)
	}
	if !strings.Contains(err.Error(), substr) {
		t.Errorf("error %q does not contain %q", err.Error(), substr)
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name       string
		password   string
		wantErr    bool
		errContain string
	}{
		{
			name:     "valid strong password",
			password: "GopherSecure!123",
			wantErr:  false,
		},
		{
			name:       "too short",
			password:   "Go!1",
			wantErr:    true,
			errContain: "password too short",
		},
		{
			name:       "missing digit",
			password:   "GopherSecure!@",
			wantErr:    true,
			errContain: "must contain at least one digit",
		},
		{
			name:       "missing special char",
			password:   "GopherSecure123",
			wantErr:    true,
			errContain: "must contain at least one special character",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePassword(tc.password)
			if tc.wantErr {
				assertErrorSubstr(t, err, tc.errContain)
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
