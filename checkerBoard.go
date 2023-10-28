package main

import (
  "fmt"
)

const X_STATE = "X"
const ZERO_STATE = "0"

func changeCurrentItem(currentItem string) string {
	if (currentItem == X_STATE) {
			currentItem = ZERO_STATE
	} else {
			currentItem = X_STATE
	}
	return currentItem;
}

func main() {
	var  columns, rows, column_width, row_height int = 5, 4, 3, 2
	var column_width_stepper, row_height_stepper int = 1, 1
	var currentItem string = "X"
	var columnsAmount = columns * column_width
	var rowsAmount = rows * row_height
	for i:=0;i<rowsAmount;i++ {
		for j:=0;j<columnsAmount;j++ {
				fmt.Print(currentItem)
				if (column_width_stepper < column_width) {
					column_width_stepper++
				} else {
					column_width_stepper = 0
					currentItem = changeCurrentItem(currentItem)
				}
		} 

		if row_height_stepper < row_height {
			row_height_stepper++
		} else {
				row_height_stepper = 0
				currentItem = changeCurrentItem(currentItem)
		}
		fmt.Println()
	}
}
