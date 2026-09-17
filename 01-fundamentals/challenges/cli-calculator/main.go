package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

var (
	ErrUsage          = errors.New("usage: calc <operation> <num1> <num2>\noperations: add, sub, mul, div, mod, pow")
	ErrDivisionByZero = errors.New("division by zero is undefined")
	ErrUnknownOp      = errors.New("unsupported operation")
)

// Calculate processes the operation on op1 and op2.
func Calculate(op string, a, b float64) (float64, error) {
	switch op {
	case "add":
		return a + b, nil
	case "sub":
		return a - b, nil
	case "mul":
		return a * b, nil
	case "div":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	case "mod":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return float64(int(a) % int(b)), nil
	case "pow":
		return math.Pow(a, b), nil
	default:
		return 0, fmt.Errorf("%w: %s", ErrUnknownOp, op)
	}
}

func run(args []string) (string, error) {
	if len(args) != 3 {
		return "", ErrUsage
	}

	op := args[0]
	a, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return "", fmt.Errorf("invalid operand %q: %w", args[1], err)
	}

	b, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return "", fmt.Errorf("invalid operand %q: %w", args[2], err)
	}

	res, err := Calculate(op, a, b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Result: %.4f", res), nil
}

func main() {
	out, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(out)
}
