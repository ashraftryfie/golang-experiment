package tempconv

import (
	"errors"
	"math"
	"testing"
)

const tolerance = 1e-4

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= tolerance
}

func TestSolutionCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name      string
		celsius   float64
		wantFahr  float64
		wantError error
	}{
		{name: "freezing point", celsius: 0.0, wantFahr: 32.0, wantError: nil},
		{name: "boiling point", celsius: 100.0, wantFahr: 212.0, wantError: nil},
		{name: "absolute zero", celsius: -273.15, wantFahr: -459.67, wantError: nil},
		{name: "below absolute zero error", celsius: -300.0, wantFahr: 0, wantError: ErrBelowAbsoluteZero},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CelsiusToFahrenheit(tc.celsius)
			if tc.wantError != nil {
				if !errors.Is(err, tc.wantError) {
					t.Fatalf("expected error %v, got %v", tc.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(got, tc.wantFahr) {
				t.Errorf("got %v, want %v", got, tc.wantFahr)
			}
		})
	}
}

func TestSolutionCelsiusToKelvin(t *testing.T) {
	got, err := CelsiusToKelvin(25.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !almostEqual(got, 298.15) {
		t.Errorf("got %v, want 298.15", got)
	}
}
