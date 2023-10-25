package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const MIN_AMOUNT_OF_ELEMENTS = 2
const AMOUNT_OF_NUMBERS_IN_RESULT = 1

// calculateRPN calculates the result of a reverse Polish notation expression.
//
// It takes a string representation of the expression as input and returns the
// calculated result as a float64 and an error if any.
func calculateRPN(s string) (float64, error) {
	var stack []float64

	// Split the input string into tokens
	var tokens = strings.Fields(s)

	// Iterate over each token
	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/", "^":
			// Check if there are enough elements on the stack to perform the operation
			if len(stack) < MIN_AMOUNT_OF_ELEMENTS {
				return 0, errors.New("not enough elements to calculate")
			}

			// Pop the top two elements from the stack
			x, y := stack[len(stack)-2], stack[len(stack)-1]
			stack = stack[:len(stack)-2]

			// Perform the operation based on the token
			switch token {
			case "+":
				stack = append(stack, x+y)
			case "-":
				stack = append(stack, x-y)
			case "*":
				stack = append(stack, x*y)
			case "/":
				if y == 0.0 {
					return 0, errors.New("can not divide by zero")
				}
				stack = append(stack, x/y)
			case "^":
				stack = append(stack, math.Pow(x, y))
			}
		default:
			// Parse the token as a float64 and append it to the stack
			res, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("token is invalid: %s", token)
			}

			stack = append(stack, res)
		}
	}

	// Check if there are only one element as the result number
	if len(stack) != AMOUNT_OF_NUMBERS_IN_RESULT {
		return 0, fmt.Errorf("incomplete expression: %.1f", stack)
	}

	return stack[0], nil
}

func main() {
	expr := "10 3 +"
	res, err := calculateRPN(expr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Result: %f\n", res)
}
