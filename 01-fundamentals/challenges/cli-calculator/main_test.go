package main

import (
	"errors"
	"math"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-4
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		op        string
		a         float64
		b         float64
		want      float64
		wantErrIs error
	}{
		{name: "addition", op: "add", a: 10, b: 5, want: 15, wantErrIs: nil},
		{name: "subtraction", op: "sub", a: 10, b: 4, want: 6, wantErrIs: nil},
		{name: "multiplication", op: "mul", a: 3, b: 7, want: 21, wantErrIs: nil},
		{name: "division", op: "div", a: 20, b: 4, want: 5, wantErrIs: nil},
		{name: "div by zero", op: "div", a: 10, b: 0, want: 0, wantErrIs: ErrDivisionByZero},
		{name: "modulo", op: "mod", a: 17, b: 5, want: 2, wantErrIs: nil},
		{name: "power", op: "pow", a: 2, b: 8, want: 256, wantErrIs: nil},
		{name: "unknown op", op: "sqrt", a: 9, b: 0, want: 0, wantErrIs: ErrUnknownOp},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Calculate(tc.op, tc.a, tc.b)
			if tc.wantErrIs != nil {
				if !errors.Is(err, tc.wantErrIs) {
					t.Fatalf("Calculate() error = %v, wantErrIs %v", err, tc.wantErrIs)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !almostEqual(got, tc.want) {
				t.Errorf("Calculate() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRunArgs(t *testing.T) {
	// Missing args
	_, err := run([]string{"add", "1"})
	if !errors.Is(err, ErrUsage) {
		t.Errorf("expected ErrUsage, got %v", err)
	}

	// Valid args
	out, err := run([]string{"mul", "4", "2.5"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "Result: 10.0000" {
		t.Errorf("got %q, want 'Result: 10.0000'", out)
	}
}
