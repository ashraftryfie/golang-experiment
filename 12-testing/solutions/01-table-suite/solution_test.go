package slug_solution

import (
	"testing"
)

func TestSlugifySolution(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "simple words", input: "Hello World", want: "hello-world"},
		{name: "with symbols", input: "Hello, World! 2024", want: "hello-world-2024"},
		{name: "multiple underscores and spaces", input: "Go_is__super   fast", want: "go-is-super-fast"},
		{name: "leading and trailing hyphens", input: "---clean-slug---", want: "clean-slug"},
		{name: "already clean slug", input: "already-clean", want: "already-clean"},
		{name: "empty string", input: "", want: ""},
		{name: "symbols only", input: "@#$%^&*!", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Slugify(tc.input)
			if got != tc.want {
				t.Errorf("Slugify(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
