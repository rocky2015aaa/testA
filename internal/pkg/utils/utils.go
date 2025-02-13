package utils

import (
	"fmt"
	"regexp"
)

func SplitByMoreThanNSpaces(text string, n int) []string {
	// Build the regular expression dynamically based on the input number of spaces
	pattern := fmt.Sprintf(`\s{%d,}`, n)
	re, _ := regexp.Compile(pattern)

	// Split the text by the regular expression
	parts := re.Split(text, -1)

	return parts
}
