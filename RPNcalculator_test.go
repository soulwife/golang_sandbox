package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

const INCOMPLETE_EXPRESSION_ERROR = "incomplete expression"
const NOT_ENOUGH_ELEMENTS_ERROR = "not enough elements to calculate"
const CAN_NOT_DIVIDE_BY_ZERO = "can not divide by zero"
const INVALID_TOKEN = "token is invalid"

// withinTolerance checks if two float64 values are within a given tolerance.
//
// It takes three parameters:
// - a: the first float64 value to compare.
// - b: the second float64 value to compare.
// - tolerance: the maximum allowed difference between a and b.
//
// It returns a boolean value indicating whether a and b are within the given tolerance.
func withinTolerance(a, b, tolerance float64) bool {
	if a == b {
		return true
	}

	absDiff := math.Abs(a - b)
	if b == 0 {
		return absDiff < tolerance
	}

	return absDiff/math.Abs(b) < tolerance
}

// All tests:
// Required input param: t: A testing.T instance used for reporting test failures.
// Returns nothing.

// Negative tests:

// Test for an empty string error
func TestCalculateRPNEmpty(t *testing.T) {
	res, err := calculateRPN("")
	expectedErrMsg := INCOMPLETE_EXPRESSION_ERROR
	if !strings.Contains(err.Error(), expectedErrMsg) || err == nil {
		t.Fatalf("Expected %s message, got %.2f, %v error", expectedErrMsg, res, err)
	}
}

// Test for error when a non-number is encountered in the RPN string
func TestCalculateRPNNotANumber(t *testing.T) {
	res, err := calculateRPN("abc")
	expectedError := fmt.Sprintf("Expected %s message and got ", INVALID_TOKEN)
	if !strings.Contains(err.Error(), INVALID_TOKEN) || err == nil {
		t.Fatalf("%s %.2f, %v error", expectedError, res, err)
	}
}

// Test checks for an "not enough elements" error
func TestCalculateRPNNotEnoughElements(t *testing.T) {
	res, err := calculateRPN("10 *")
	expectedError := NOT_ENOUGH_ELEMENTS_ERROR
	if err.Error() != expectedError || err == nil {
		t.Fatalf("Expected %s message and got %.2f, %v error", expectedError, res, err)
	}
}

// Test for "not enough elements" error in case of incomplete expression that is short of operand
func TestCalculateRPNNotEnoughOperandElements(t *testing.T) {
	expr := "2 3 + 4"
	res, err := calculateRPN(expr)
	expectedErr := "incomplete expression"
	if !strings.Contains(err.Error(), expectedErr) || err == nil {
		t.Fatalf("Expected %s message and got %.2f, %v error",
			expectedErr, res, err)
	}
}

// test checks for a "divide by zero" error
func TestCalculateRPNDivideByZero(t *testing.T) {
	res, err := calculateRPN("10 0 /")
	expectedError := "can not divide by zero"
	if err.Error() != expectedError || err == nil {
		t.Fatalf("Expected %s message and got %.2f, %v error", expectedError, res, err)
	}
}

// Positive tests:

// Test for calculating RPN with one number
func TestCalculateRPNOneNumber(t *testing.T) {
	num := 10.0
	input := strconv.FormatFloat(num, 'f', 0, 64)
	res, err := calculateRPN(input)
	if err != nil {
		t.Errorf(`Expected %.2f, got %.2f`, res, num)
	}
}

// Tests the calculateRPN function with various expressions
// It verifies that the function returns the expected result
// within a certain tolerance.
func TestCalculateRPNWithNumbers(t *testing.T) {
	testCases := []struct {
		expression string
		expected   float64
	}{
		{"2 3 +", 5},
		{"2 3 -", -1},
		{"2 3 *", 6},
		{"2 3 /", 0.66},
		{"2 3 ^", 8},
		{"2.5 3.05 +", 5.55},
		{"10 3 2 + -", 5},
		{"10 3 * 2 ^", 900},
		{"1 2 3 4 5 6 7 8 9 10 + + + + + + + + +", 55},
		{"1 sin", 0.841},
		{"180 cos", -0.598},
		{"360 tan", -3.383},
		{"90 ctg", -0.501},
		{"-1 acos", 3.141},
		{"1 asin", 1.57},
		{"1 atan", 0.785},
		{"4 sqrt", 2},
		{"4 16 sqrt sqrt +", 6},
		{"4 16 16 sqrt sqrt + +", 22},
	}

	for _, tc := range testCases {
		res, err := calculateRPN(tc.expression)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !withinTolerance(res, tc.expected, 1e-1) {
			t.Errorf("Expected %v, got %v", tc.expected, res)
		}
	}
}
