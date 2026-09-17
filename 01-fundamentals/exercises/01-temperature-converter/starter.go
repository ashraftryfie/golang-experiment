package tempconv

import (
	"errors"
)

// ErrNotImplemented is returned by starter templates before you implement them.
var ErrNotImplemented = errors.New("TODO: implement function")

// ErrBelowAbsoluteZero is returned when an input temperature is physically impossible.
var ErrBelowAbsoluteZero = errors.New("temperature is below absolute zero")

const AbsoluteZeroC = -273.15

// CelsiusToFahrenheit converts Celsius to Fahrenheit.
// Returns ErrBelowAbsoluteZero if c < -273.15.
func CelsiusToFahrenheit(c float64) (float64, error) {
	// TODO: Replace with your implementation
	return 0, ErrNotImplemented
}

// FahrenheitToCelsius converts Fahrenheit to Celsius.
// Returns ErrBelowAbsoluteZero if resulting Celsius is below AbsoluteZeroC.
func FahrenheitToCelsius(f float64) (float64, error) {
	// TODO: Replace with your implementation
	return 0, ErrNotImplemented
}

// CelsiusToKelvin converts Celsius to Kelvin.
// Returns ErrBelowAbsoluteZero if c < AbsoluteZeroC.
func CelsiusToKelvin(c float64) (float64, error) {
	// TODO: Replace with your implementation
	return 0, ErrNotImplemented
}
