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

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name      string
		celsius   float64
		wantFahr  float64
		wantError error
	}{
		{name: "freezing point of water", celsius: 0.0, wantFahr: 32.0, wantError: nil},
		{name: "boiling point of water", celsius: 100.0, wantFahr: 212.0, wantError: nil},
		{name: "body temperature", celsius: 37.0, wantFahr: 98.6, wantError: nil},
		{name: "absolute zero", celsius: -273.15, wantFahr: -459.67, wantError: nil},
		{name: "below absolute zero error", celsius: -300.0, wantFahr: 0, wantError: ErrBelowAbsoluteZero},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CelsiusToFahrenheit(tc.celsius)
			if errors.Is(err, ErrNotImplemented) {
				t.Skip("skipping: CelsiusToFahrenheit is not yet implemented (implement in starter.go)")
			}
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
				t.Errorf("CelsiusToFahrenheit(%v) = %v; want %v", tc.celsius, got, tc.wantFahr)
			}
		})
	}
}

func TestFahrenheitToCelsius(t *testing.T) {
	tests := []struct {
		name      string
		fahr      float64
		wantC     float64
		wantError error
	}{
		{name: "freezing point", fahr: 32.0, wantC: 0.0, wantError: nil},
		{name: "boiling point", fahr: 212.0, wantC: 100.0, wantError: nil},
		{name: "below absolute zero error", fahr: -500.0, wantC: 0, wantError: ErrBelowAbsoluteZero},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FahrenheitToCelsius(tc.fahr)
			if errors.Is(err, ErrNotImplemented) {
				t.Skip("skipping: FahrenheitToCelsius is not yet implemented (implement in starter.go)")
			}
			if tc.wantError != nil {
				if !errors.Is(err, tc.wantError) {
					t.Fatalf("expected error %v, got %v", tc.wantError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(got, tc.wantC) {
				t.Errorf("FahrenheitToCelsius(%v) = %v; want %v", tc.fahr, got, tc.wantC)
			}
		})
	}
}
