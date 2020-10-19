package utils

import "fmt"

// Percentage returns a string with the representation of the percentage of X in Y
func Percentage(x, y float64) string {
	if y == 0 {
		return "0%"
	}

	return fmt.Sprintf("%.2f%%", (x/y)*100)
}
