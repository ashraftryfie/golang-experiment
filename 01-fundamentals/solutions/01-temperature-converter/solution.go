package tempconv

import (
	"errors"
)

// ErrBelowAbsoluteZero is returned when an input temperature is physically impossible.
var ErrBelowAbsoluteZero = errors.New("temperature is below absolute zero")

const AbsoluteZeroC = -273.15

// CelsiusToFahrenheit converts Celsius to Fahrenheit.
func CelsiusToFahrenheit(c float64) (float64, error) {
	if c < AbsoluteZeroC {
		return 0, ErrBelowAbsoluteZero
	}
	return (c * 9.0 / 5.0) + 32.0, nil
}

// FahrenheitToCelsius converts Fahrenheit to Celsius.
func FahrenheitToCelsius(f float64) (float64, error) {
	c := (f - 32.0) * 5.0 / 9.0
	if c < AbsoluteZeroC {
		return 0, ErrBelowAbsoluteZero
	}
	return c, nil
}

// CelsiusToKelvin converts Celsius to Kelvin.
func CelsiusToKelvin(c float64) (float64, error) {
	if c < AbsoluteZeroC {
		return 0, ErrBelowAbsoluteZero
	}
	return c + 273.15, nil
}
