package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const ONE_ELEMENT = 1
const TWO_ELEMENTS = 2

// checkAmountOfElements checks if the given amount of elements is less than the expected amount.
//
// It takes in two parameters:
// - givenAmountOfElements: the actual amount of elements
// - expectedAmount: the expected amount of elements
//
// It returns an error if the given amount of elements is less than the expected amount.
func checkAmountOfElements(givenAmountOfElements int, expectedAmount int) error {
	if givenAmountOfElements < expectedAmount {
		return errors.New("not enough elements to calculate")
	}

	return nil
}

// calculateRPN calculates the result of a reverse Polish notation expression.
//
// It takes a string representation of the expression as input and returns the
// calculated result as a float64 and an error if any.
func calculateRPN(s string) (float64, error) {
	var stack []float64

	// Split the input string into tokens
	var tokens = strings.Fields(s)

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/", "^":
			err := checkAmountOfElements(len(stack), TWO_ELEMENTS)
			if err != nil {
				return 0, err
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
				if y == 0 {
					return 0, errors.New("can not divide by zero")
				}
				stack = append(stack, x/y)
			case "^":
				stack = append(stack, math.Pow(x, y))
			}
		case "sin", "cos", "tan", "asin", "acos", "atan", "sqrt", "ctg":
			err := checkAmountOfElements(len(stack), ONE_ELEMENT)
			if err != nil {
				return 0, err
			}

			x := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			switch token {
			case "sin":
				stack = append(stack, math.Sin(x))
			case "cos":
				stack = append(stack, math.Cos(x))
			case "tan":
				stack = append(stack, math.Tan(x))
			case "asin":
				stack = append(stack, math.Asin(x))
			case "acos":
				stack = append(stack, math.Acos(x))
			case "atan":
				stack = append(stack, math.Atan(x))
			case "ctg":
				stack = append(stack, 1/math.Tan(x))
			case "sqrt":
				stack = append(stack, math.Sqrt(x))
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
	if len(stack) != ONE_ELEMENT {
		return 0, fmt.Errorf("incomplete expression: %.1f", stack)
	}

	return stack[0], nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter an expression:")

	for scanner.Scan() {
		result, err := calculateRPN(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}

		formattedResult := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.3f", result), "0"), ".")
		fmt.Println("Result: " + formattedResult)
		fmt.Println("Enter an expression:")
	}

}
