package slug

import (
	"testing"
)

func TestSlugify(t *testing.T) {
	// Table-driven test suite
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// TODO: Complete test cases to verify Slugify edge cases
		{name: "TODO: implement test cases", input: "TODO", want: "TODO"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "TODO: implement test cases" {
				t.Skip("skipping: table test cases not yet implemented (implement in starter_test.go)")
			}
			got := Slugify(tc.input)
			if got != tc.want {
				t.Errorf("Slugify(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
