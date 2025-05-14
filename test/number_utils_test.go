package test

import (
	"gfly/app/utils"
	"testing"
)

func TestCeiling(t *testing.T) {
	tests := []struct {
		name      string
		number    float64
		precision []int
		expected  float64
	}{
		{"Integer", 5.0, nil, 5.0},
		{"RoundUp", 5.1, nil, 6.0},
		{"RoundUpWithPrecision", 5.123, []int{2}, 5.13},
		{"NegativeNumber", -5.1, nil, -5.0},
		{"NegativeWithPrecision", -5.123, []int{2}, -5.12},
		{"ZeroPrecision", 5.123, []int{0}, 6.0},
		{"LargeNumber", 9999.9999, []int{2}, 10000.00},
		{"MultiplePrecisionArgs", 5.123, []int{2, 3}, 5.13}, // Should use only the first precision value
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumCeiling(test.number, test.precision...)
			if result != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, result)
			}
		})
	}
}

func TestFloor(t *testing.T) {
	tests := []struct {
		name      string
		number    float64
		precision []int
		expected  float64
	}{
		{"Integer", 5.0, nil, 5.0},
		{"RoundDown", 5.9, nil, 5.0},
		{"RoundDownWithPrecision", 5.129, []int{2}, 5.12},
		{"NegativeNumber", -5.1, nil, -6.0},
		{"NegativeWithPrecision", -5.129, []int{2}, -5.13},
		{"ZeroPrecision", 5.9, []int{0}, 5.0},
		{"LargeNumber", 10000.0001, []int{2}, 10000.00},
		{"MultiplePrecisionArgs", 5.129, []int{2, 3}, 5.12}, // Should use only the first precision value
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumFloor(test.number, test.precision...)
			if result != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, result)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name               string
		number             float64
		decimals           int
		decimalSeparator   string
		thousandsSeparator string
		expected           string
	}{
		{"NoDecimals", 1234.56, 0, ".", ",", "1,234"},
		{"TwoDecimals", 1234.56, 2, ".", ",", "1,234.56"},
		{"CustomSeparators", 1234.56, 2, ",", " ", "1 234,56"},
		{"LargeNumber", 1234567.89, 2, ".", ",", "1,234,567.89"},
		{"NegativeNumber", -1234.56, 2, ".", ",", "-1,234.56"},
		{"ZeroDecimals", 0.0, 0, ".", ",", "0"},
		{"EmptySeparators", 1234.56, 2, ".", "", "1234.56"},
		{"MoreDecimalsThanProvided", 1234.5, 3, ".", ",", "1,234.500"},
		{"VeryLargeNumber", 1234567890.12, 2, ".", ",", "1,234,567,890.12"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumFormat(test.number, test.decimals, test.decimalSeparator, test.thousandsSeparator)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestFormatCompact(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		decimals int
		expected string
	}{
		{"LessThanThousand", 999, 0, "999"},
		{"Thousand", 1000, 0, "1K"},
		{"ThousandWithDecimals", 1500, 1, "1.5K"},
		{"Million", 1000000, 0, "1M"},
		{"MillionWithDecimals", 1500000, 1, "1.5M"},
		{"Billion", 1000000000, 0, "1B"},
		{"BillionWithDecimals", 1500000000, 1, "1.5B"},
		{"NegativeNumber", -1500, 1, "-1.5K"},
		{"Zero", 0, 2, "0.00"},
		{"SmallNumber", 0.123, 2, "0.12"},
		{"ExactThousandBoundary", 999.9, 1, "999.9"},
		{"HighPrecision", 1234567.89, 3, "1.235M"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumFormatCompact(test.number, test.decimals)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestFormatPercentage(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		decimals int
		expected string
	}{
		{"ZeroPercent", 0, 0, "0%"},
		{"HundredPercent", 1, 0, "100%"},
		{"FiftyPercent", 0.5, 0, "50%"},
		{"WithDecimals", 0.1234, 2, "12.34%"},
		{"NegativePercent", -0.5, 0, "-50%"},
		{"HighPrecision", 0.12345, 4, "12.3450%"},
		{"VerySmallNumber", 0.000123, 5, "0.01230%"},
		{"LargeNumber", 12.34, 2, "1234.00%"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumFormatPercentage(test.number, test.decimals)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestIsEven(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		expected bool
	}{
		{"Zero", 0, true},
		{"PositiveEven", 2, true},
		{"PositiveOdd", 3, false},
		{"NegativeEven", -4, true},
		{"NegativeOdd", -5, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumIsEven(test.number)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestIsOdd(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		expected bool
	}{
		{"Zero", 0, false},
		{"PositiveEven", 2, false},
		{"PositiveOdd", 3, true},
		{"NegativeEven", -4, false},
		{"NegativeOdd", -5, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumIsOdd(test.number)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		name     string
		number   float64
		total    float64
		decimals []int
		expected float64
	}{
		{"ZeroPercent", 0, 100, nil, 0},
		{"HundredPercent", 100, 100, nil, 100},
		{"FiftyPercent", 50, 100, nil, 50},
		{"WithDecimals", 12.34, 100, []int{2}, 12.34},
		{"RoundedDecimals", 12.345, 100, []int{2}, 12.35},
		{"ZeroTotal", 50, 0, nil, 0},
		{"NegativeNumber", -25, 100, nil, -25},
		{"NegativeTotal", 25, -100, nil, -25},
		{"BothNegative", -25, -100, nil, 25},
		{"HighPrecision", 1, 3, []int{5}, 33.33333},
		{"MultiplePrecisionArgs", 12.345, 100, []int{2, 3}, 12.35}, // Should use only the first precision value
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.NumPercent(test.number, test.total, test.decimals...)
			if result != test.expected {
				t.Errorf("Expected %f, got %f", test.expected, result)
			}
		})
	}
}
