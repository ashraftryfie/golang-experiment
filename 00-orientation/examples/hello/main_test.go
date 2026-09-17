package main

import "testing"

func TestGreet(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid name provided",
			input:    "Ashraf",
			expected: "Hello, Ashraf! Welcome to Go.",
		},
		{
			name:     "empty string defaults to Gopher",
			input:    "",
			expected: "Hello, Gopher! Welcome to Go.",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Greet(tc.input)
			if result != tc.expected {
				t.Errorf("Greet(%q) = %q; want %q", tc.input, result, tc.expected)
			}
		})
	}
}
