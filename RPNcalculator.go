package main

import (
  "fmt"
  "strconv"
  "strings"
	"errors"
	"math"
)

const MIN_AMOUNT_OF_ELEMENTS = 2
const RESULT_NUMBERS_AMOUNT = 1

func calculateRPN(s string) (float64, error) {
  var stack []float64
  var tokens = strings.Fields(s)
  var result float64

  for _, token := range tokens {
    switch token {
    case "+", "-", "*", "/", "^":
      var x, y float64
      if len(stack) < MIN_AMOUNT_OF_ELEMENTS {
        return 0, errors.New("There are not enough elemenents to calculate")
      }
      stack, x, y = stack[:len(stack)-2], stack[len(stack)-2:][0], stack[len(stack)-2:][1]

      switch token {
				case "+":
					result = x + y
				case "-":
					result = x - y
				case "*":
					result = x * y
				case "/":
					if y == 0.0 {
						return 0, errors.New("We can not divide by zero")
					}
					result = x / y
				case "^":
          result = math.Pow(x, y)
      }
    default:
			var err error
      result, err = strconv.ParseFloat(token, 64)
      if err != nil {
        return 0, errors.New(fmt.Sprintf("Token is invalid: %s", token))
      }
    }
    stack = append(stack, result)
  }

  if len(stack) != RESULT_NUMBERS_AMOUNT {
    return 0, errors.New(fmt.Sprintf("There are incomplete expression: %.1f", stack))
  }

  return stack[0], nil
}

func main() {
	expression := "10 -5 *"
  result, err := calculateRPN(expression)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("Result: " + strconv.FormatFloat(result, 'f', -1, 64))
  }
}
