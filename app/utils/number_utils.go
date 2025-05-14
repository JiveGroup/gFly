package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Ceiling rounds a number up to the nearest integer or precision.
// If precision is provided, the number is rounded up to that decimal place.
func NumCeiling(number float64, precision ...int) float64 {
	if len(precision) > 0 {
		factor := math.Pow(10, float64(precision[0]))
		return math.Ceil(number*factor) / factor
	}
	return math.Ceil(number)
}

// Floor rounds a number down to the nearest integer or precision.
// If precision is provided, the number is rounded down to that decimal place.
func NumFloor(number float64, precision ...int) float64 {
	if len(precision) > 0 {
		factor := math.Pow(10, float64(precision[0]))
		return math.Floor(number*factor) / factor
	}
	return math.Floor(number)
}

// Format formats a number with grouped thousands.
// The default separator is a comma, but a custom separator can be provided.
func NumFormat(number float64, decimals int, decimalSeparator, thousandsSeparator string) string {
	// If decimals is 0, truncate the number instead of rounding
	if decimals == 0 {
		number = math.Floor(number)
	}

	// Format the number with the specified number of decimal places
	formatStr := "%." + strconv.Itoa(decimals) + "f"
	formattedNumber := fmt.Sprintf(formatStr, number)

	// Split the number into integer and decimal parts
	parts := strings.Split(formattedNumber, ".")
	integerPart := parts[0]

	// Add thousands separator
	var result strings.Builder
	for i, char := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			result.WriteString(thousandsSeparator)
		}
		result.WriteRune(char)
	}

	// Add decimal part if needed
	if decimals > 0 {
		result.WriteString(decimalSeparator)
		if len(parts) > 1 {
			result.WriteString(parts[1])
		} else {
			result.WriteString(strings.Repeat("0", decimals))
		}
	}

	return result.String()
}

// FormatCompact formats a number to a compact form (e.g., 1K, 1M).
func NumFormatCompact(number float64, decimals int) string {
	absNumber := math.Abs(number)
	sign := ""
	if number < 0 {
		sign = "-"
	}

	switch {
	case absNumber >= 1_000_000_000:
		return sign + fmt.Sprintf("%.*f", decimals, absNumber/1_000_000_000) + "B"
	case absNumber >= 1_000_000:
		return sign + fmt.Sprintf("%.*f", decimals, absNumber/1_000_000) + "M"
	case absNumber >= 1_000:
		return sign + fmt.Sprintf("%.*f", decimals, absNumber/1_000) + "K"
	default:
		return sign + fmt.Sprintf("%.*f", decimals, absNumber)
	}
}

// FormatPercentage formats a number as a percentage.
func NumFormatPercentage(number float64, decimals int) string {
	return fmt.Sprintf("%.*f%%", decimals, number*100)
}

// IsEven determines if a number is even.
func NumIsEven(number int) bool {
	return number%2 == 0
}

// NumIsOdd determines if a number is odd.
func NumIsOdd(number int) bool {
	return number%2 != 0
}

// NumPercent calculates the percentage of a number.
func NumPercent(number, total float64, decimals ...int) float64 {
	if total == 0 {
		return 0
	}

	percentage := (number / total) * 100

	if len(decimals) > 0 {
		factor := math.Pow(10, float64(decimals[0]))
		return math.Round(percentage*factor) / factor
	}

	return percentage
}
