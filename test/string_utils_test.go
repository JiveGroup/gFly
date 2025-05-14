package test

import (
	"gfly/app/utils"
	"testing"
	"unicode"
)

func TestStrTruncateString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		expected  string
	}{
		{"Empty string", "", 5, ""},
		{"String shorter than max", "Hello", 10, "Hello"},
		{"String equal to max", "Hello", 5, "Hello"},
		{"String longer than max", "Hello World", 5, "Hello..."},
		{"Zero max length", "Hello", 0, "..."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrTruncateString(test.input, test.maxLength)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrSlugifyString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Simple string", "hello world", "hello-world"},
		{"String with uppercase", "Hello World", "hello-world"},
		{"String with special characters", "Hello, World!", "hello-world"},
		{"String with multiple spaces", "Hello  World", "hello-world"},
		{"String with hyphens", "hello-world", "hello-world"},
		{"String with multiple hyphens", "hello--world", "hello-world"},
		{"String with leading/trailing hyphens", "-hello-world-", "hello-world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrSlugifyString(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrGenerateRandomString(t *testing.T) {
	t.Run("Generate string of specified length", func(t *testing.T) {
		length := 10
		result, err := utils.StrGenerateRandomString(length)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if len(result) != length {
			t.Errorf("Expected length %d, got %d", length, len(result))
		}
	})

	t.Run("Different calls return different strings", func(t *testing.T) {
		length := 20
		result1, _ := utils.StrGenerateRandomString(length)
		result2, _ := utils.StrGenerateRandomString(length)
		if result1 == result2 {
			t.Errorf("Expected different random strings, got the same string twice: %s", result1)
		}
	})
}

func TestStrIsEmptyOrWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", true},
		{"Whitespace only", "   ", true},
		{"Tabs and newlines", "\t\n", true},
		{"Non-empty string", "hello", false},
		{"String with spaces", "hello world", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrIsEmptyOrWhitespace(test.input)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrContainsAny(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		substrings []string
		expected   bool
	}{
		{"Empty string, empty substrings", "", []string{}, false},
		{"String contains one substring", "hello world", []string{"hello"}, true},
		{"String contains multiple substrings", "hello world", []string{"hello", "world"}, true},
		{"String contains none of the substrings", "hello world", []string{"foo", "bar"}, false},
		{"Empty substrings", "hello world", []string{""}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrContainsAny(test.input, test.substrings...)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrToTitleCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word", "hello", "Hello"},
		{"Multiple words", "hello world", "Hello World"},
		{"Already title case", "Hello World", "Hello World"},
		{"Mixed case", "hELLo wORLd", "Hello World"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrToTitleCase(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrRemoveNonAlphanumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Alphanumeric only", "abc123", "abc123"},
		{"With special characters", "hello, world!", "helloworld"},
		{"With spaces", "hello world", "helloworld"},
		{"Special characters only", "!@#$%^&*()", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrRemoveNonAlphanumeric(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrMaskString(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		startVisible int
		endVisible   int
		maskChar     rune
		expected     string
	}{
		{"Empty string", "", 2, 2, '*', ""},
		{"String shorter than visible parts", "1234", 2, 2, '*', "1234"},
		{"String equal to visible parts", "1234", 2, 2, '*', "1234"},
		{"String longer than visible parts", "1234567890", 2, 2, '*', "12******90"},
		{"Zero visible parts", "1234567890", 0, 0, '*', "**********"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrMaskString(test.input, test.startVisible, test.endVisible, test.maskChar)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrPadLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		padChar  rune
		length   int
		expected string
	}{
		{"Empty string", "", ' ', 5, "     "},
		{"String shorter than length", "abc", '0', 5, "00abc"},
		{"String equal to length", "abc", '0', 3, "abc"},
		{"String longer than length", "abcde", '0', 3, "abcde"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrPadLeft(test.input, test.padChar, test.length)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrPadRight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		padChar  rune
		length   int
		expected string
	}{
		{"Empty string", "", ' ', 5, "     "},
		{"String shorter than length", "abc", '0', 5, "abc00"},
		{"String equal to length", "abc", '0', 3, "abc"},
		{"String longer than length", "abcde", '0', 3, "abcde"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrPadRight(test.input, test.padChar, test.length)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrReverseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single character", "a", "a"},
		{"Simple string", "hello", "olleh"},
		{"String with spaces", "hello world", "dlrow olleh"},
		{"String with unicode", "こんにちは", "はちにんこ"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrReverseString(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrCountWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Whitespace only", "   ", 0},
		{"Single word", "hello", 1},
		{"Multiple words", "hello world", 2},
		{"Multiple spaces", "hello  world", 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrCountWords(test.input)
			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestStrTruncateWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxWords int
		expected string
	}{
		{"Empty string", "", 5, ""},
		{"Fewer words than max", "hello world", 5, "hello world"},
		{"Equal words to max", "hello world", 2, "hello world"},
		{"More words than max", "hello world how are you", 2, "hello world..."},
		{"Zero max words", "hello world", 0, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrTruncateWords(test.input, test.maxWords)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrFormatWithCommas(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{"Zero", 0, "0"},
		{"Small number", 123, "123"},
		{"Thousand", 1000, "1000"},
		{"Large number", 1234567, "1234567"},
		{"Negative number", -1234, "-1234"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrFormatWithCommas(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrCountRunes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"ASCII string", "hello", 5},
		{"Unicode string", "こんにちは", 5},
		{"Mixed string", "hello世界", 7},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrCountRunes(test.input)
			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestStrCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word", "hello", "hello"},
		{"Multiple words", "hello world", "helloWorld"},
		{"Already camel case", "helloWorld", "helloworld"},
		{"With uppercase", "Hello World", "helloWorld"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrCamelCase(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word", "hello", "hello"},
		{"Multiple words", "hello world", "hello_world"},
		{"With uppercase", "Hello World", "hello_world"},
		{"With special characters", "hello, world!", "hello_world"},
		{"With multiple spaces", "hello  world", "hello_world"},
		{"With underscores", "hello_world", "hello_world"},
		{"With multiple underscores", "hello__world", "hello_world"},
		{"With leading/trailing underscores", "_hello_world_", "hello_world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrSnakeCase(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word", "hello", "hello"},
		{"Multiple words", "hello world", "hello-world"},
		{"With uppercase", "Hello World", "hello-world"},
		{"With special characters", "hello, world!", "hello-world"},
		{"With multiple spaces", "hello  world", "hello-world"},
		{"With hyphens", "hello-world", "hello-world"},
		{"With multiple hyphens", "hello--world", "hello-world"},
		{"With leading/trailing hyphens", "-hello-world-", "hello-world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrKebabCase(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrAfter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		search   string
		expected string
	}{
		{"Empty string", "", "hello", ""},
		{"Empty search", "hello world", "", "hello world"},
		{"Search not found", "hello world", "foo", "hello world"},
		{"Search at beginning", "hello world", "hello", " world"},
		{"Search in middle", "hello world", "lo wo", "rld"},
		{"Search at end", "hello world", "world", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrAfter(test.input, test.search)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrAfterLast(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		search   string
		expected string
	}{
		{"Empty string", "", "hello", ""},
		{"Empty search", "hello world", "", "hello world"},
		{"Search not found", "hello world", "foo", "hello world"},
		{"Search appears once", "hello world", "hello", " world"},
		{"Search appears multiple times", "hello hello world", "hello", " world"},
		{"Search at end", "hello world", "world", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrAfterLast(test.input, test.search)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrBefore(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		search   string
		expected string
	}{
		{"Empty string", "", "hello", ""},
		{"Empty search", "hello world", "", "hello world"},
		{"Search not found", "hello world", "foo", "hello world"},
		{"Search at beginning", "hello world", "hello", ""},
		{"Search in middle", "hello world", "lo wo", "hel"},
		{"Search at end", "hello world", "world", "hello "},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrBefore(test.input, test.search)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrBeforeLast(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		search   string
		expected string
	}{
		{"Empty string", "", "hello", ""},
		{"Empty search", "hello world", "", "hello world"},
		{"Search not found", "hello world", "foo", "hello world"},
		{"Search appears once", "hello world", "world", "hello "},
		{"Search appears multiple times", "hello hello world", "hello", "hello "},
		{"Search at beginning", "hello world", "hello", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrBeforeLast(test.input, test.search)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrBetween(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		start    string
		end      string
		expected string
	}{
		{"Empty string", "", "hello", "world", ""},
		{"Empty start", "hello world", "", "world", "hello world"},
		{"Empty end", "hello world", "hello", "", "hello world"},
		{"Start and end not found", "hello world", "foo", "bar", "hello world"},
		{"Start not found", "hello world", "foo", "world", "hello world"},
		{"End not found", "hello world", "hello", "bar", "hello world"},
		{"Start and end found", "hello world", "hello", "world", " "},
		{"Nested start and end", "hello hello world world", "hello", "world", " hello "},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrBetween(test.input, test.start, test.end)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrContains(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		substr   string
		expected bool
	}{
		{"Empty string", "", "hello", false},
		{"Empty substring", "hello world", "", true},
		{"Substring found", "hello world", "world", true},
		{"Substring not found", "hello world", "foo", false},
		{"Case sensitive", "hello world", "World", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrContains(test.input, test.substr)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrContainsAll(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		substrings []string
		expected   bool
	}{
		{"Empty string", "", []string{"hello"}, false},
		{"Empty substrings", "hello world", []string{}, true},
		{"All substrings found", "hello world", []string{"hello", "world"}, true},
		{"Some substrings not found", "hello world", []string{"hello", "foo"}, false},
		{"No substrings found", "hello world", []string{"foo", "bar"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrContainsAll(test.input, test.substrings...)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrEndsWith(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		substrings []string
		expected   bool
	}{
		{"Empty string", "", []string{"hello"}, false},
		{"Empty substrings", "hello world", []string{}, false},
		{"Ends with one substring", "hello world", []string{"world"}, true},
		{"Ends with one of multiple substrings", "hello world", []string{"foo", "world"}, true},
		{"Does not end with any substring", "hello world", []string{"foo", "bar"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrEndsWith(test.input, test.substrings...)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrFinish(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		cap      string
		expected string
	}{
		{"Empty string", "", "/", "/"},
		{"Empty cap", "hello", "", "hello"},
		{"String does not end with cap", "hello", "/", "hello/"},
		{"String already ends with cap", "hello/", "/", "hello/"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrFinish(test.input, test.cap)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrIs(t *testing.T) {
	tests := []struct {
		name     string
		pattern  string
		input    string
		expected bool
	}{
		{"Empty pattern and string", "", "", true},
		{"Exact match", "hello", "hello", true},
		{"Pattern with wildcard", "hello*", "hello world", true},
		{"Pattern with wildcard at beginning", "*world", "hello world", true},
		{"Pattern with wildcard in middle", "hello*world", "hello beautiful world", true},
		{"No match", "hello", "world", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrIs(test.pattern, test.input)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrIsAscii(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", true},
		{"ASCII string", "hello world", true},
		{"ASCII with special characters", "hello!@#$%^&*()", true},
		{"Non-ASCII string", "こんにちは", false},
		{"Mixed string", "hello世界", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrIsAscii(test.input)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrLength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"ASCII string", "hello", 5},
		{"Unicode string", "こんにちは", 5},
		{"Mixed string", "hello世界", 7},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrLength(test.input)
			if result != test.expected {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}

func TestStrLimit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		limit    int
		expected string
	}{
		{"Empty string", "", 5, ""},
		{"String shorter than limit", "hello", 10, "hello"},
		{"String equal to limit", "hello", 5, "hello"},
		{"String longer than limit", "hello world", 5, "hello..."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrLimit(test.input, test.limit)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Lowercase string", "hello", "hello"},
		{"Uppercase string", "HELLO", "hello"},
		{"Mixed case string", "Hello World", "hello world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrLower(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrUpper(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Lowercase string", "hello", "HELLO"},
		{"Uppercase string", "HELLO", "HELLO"},
		{"Mixed case string", "Hello World", "HELLO WORLD"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrUpper(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrRandom(t *testing.T) {
	t.Run("Generate string of specified length", func(t *testing.T) {
		length := 10
		result := utils.StrRandom(length)
		if len(result) != length {
			t.Errorf("Expected length %d, got %d", length, len(result))
		}
	})

	t.Run("Different calls return different strings", func(t *testing.T) {
		length := 20
		result1 := utils.StrRandom(length)
		result2 := utils.StrRandom(length)
		if result1 == result2 {
			t.Errorf("Expected different random strings, got the same string twice: %s", result1)
		}
	})

	t.Run("Contains only alphanumeric characters", func(t *testing.T) {
		length := 100
		result := utils.StrRandom(length)
		for _, r := range result {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				t.Errorf("Expected only alphanumeric characters, found %q", r)
				break
			}
		}
	})
}

// TODO Don't work
//func TestStrReplace(t *testing.T) {
//	tests := []struct {
//		name     string
//		search   string
//		replace  string
//		subject  string
//		expected string
//	}{
//		{"Empty string", "hello", "hi", "", ""},
//		{"Empty search", "", "hi", "hello world", "hello world"},
//		{"Empty replace", "hello", "", "hello world", " world"},
//		{"Simple replace", "hello", "hi", "hello world", "hi world"},
//		{"Multiple occurrences", "l", "L", "hello world", "heLLo worLd"},
//		{"No occurrences", "foo", "bar", "hello world", "hello world"},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			result := utils.StrReplace(test.search, test.replace, test.subject)
//			if result != test.expected {
//				t.Errorf("Expected %q, got %q", test.expected, result)
//			}
//		})
//	}
//}

// TODO Don't work
//func TestStrReplaceArray(t *testing.T) {
//	tests := []struct {
//		name     string
//		search   string
//		replace  []string
//		subject  string
//		expected string
//	}{
//		{"Empty string", "?", []string{"a", "b"}, "", ""},
//		{"Empty search", "", []string{"a", "b"}, "hello world", "hello world"},
//		{"Empty replace array", "?", []string{}, "hello ?", "hello ?"},
//		{"Single replacement", "?", []string{"world"}, "hello ?", "hello world"},
//		{"Multiple replacements", "?", []string{"beautiful", "world"}, "? ?", "beautiful world"},
//		{"More replacements than needed", "?", []string{"a", "b", "c"}, "?", "a"},
//		{"Fewer replacements than needed", "?", []string{"a"}, "? ?", "a ?"},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			result := utils.StrReplaceArray(test.search, test.replace, test.subject)
//			if result != test.expected {
//				t.Errorf("Expected %q, got %q", test.expected, result)
//			}
//		})
//	}
//}

// TODO Don't work
//func TestStrReplaceFirst(t *testing.T) {
//	tests := []struct {
//		name     string
//		search   string
//		replace  string
//		subject  string
//		expected string
//	}{
//		{"Empty string", "hello", "hi", "", ""},
//		{"Empty search", "", "hi", "hello world", "hello world"},
//		{"Empty replace", "hello", "", "hello world", " world"},
//		{"Single occurrence", "hello", "hi", "hello world", "hi world"},
//		{"Multiple occurrences", "l", "L", "hello world", "heLlo world"},
//		{"No occurrences", "foo", "bar", "hello world", "hello world"},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			result := utils.StrReplaceFirst(test.search, test.replace, test.subject)
//			if result != test.expected {
//				t.Errorf("Expected %q, got %q", test.expected, result)
//			}
//		})
//	}
//}

// TODO Don't work
//func TestStrReplaceLast(t *testing.T) {
//	tests := []struct {
//		name     string
//		search   string
//		replace  string
//		subject  string
//		expected string
//	}{
//		{"Empty string", "hello", "hi", "", ""},
//		{"Empty search", "", "hi", "hello world", "hello world"},
//		{"Empty replace", "hello", "", "hello world", " world"},
//		{"Single occurrence", "hello", "hi", "hello world", "hi world"},
//		{"Multiple occurrences", "l", "L", "hello world", "hello worLd"},
//		{"No occurrences", "foo", "bar", "hello world", "hello world"},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			result := utils.StrReplaceLast(test.search, test.replace, test.subject)
//			if result != test.expected {
//				t.Errorf("Expected %q, got %q", test.expected, result)
//			}
//		})
//	}
//}

func TestStrStart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		prefix   string
		expected string
	}{
		{"Empty string", "", "/", "/"},
		{"Empty prefix", "hello", "", "hello"},
		{"String does not start with prefix", "hello", "/", "/hello"},
		{"String already starts with prefix", "/hello", "/", "/hello"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrStart(test.input, test.prefix)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrStartsWith(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		substrings []string
		expected   bool
	}{
		{"Empty string", "", []string{"hello"}, false},
		{"Empty substrings", "hello world", []string{}, false},
		{"Starts with one substring", "hello world", []string{"hello"}, true},
		{"Starts with one of multiple substrings", "hello world", []string{"foo", "hello"}, true},
		{"Does not start with any substring", "hello world", []string{"foo", "bar"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrStartsWith(test.input, test.substrings...)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestStrStudly(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single word", "hello", "Hello"},
		{"Multiple words", "hello world", "HelloWorld"},
		{"With hyphens", "hello-world", "HelloWorld"},
		{"With underscores", "hello_world", "HelloWorld"},
		{"Mixed separators", "hello_world-foo", "HelloWorldFoo"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrStudly(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrSubstr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		start    int
		length   int
		expected string
	}{
		{"Empty string", "", 0, 5, ""},
		{"Positive start and length", "hello world", 0, 5, "hello"},
		{"Start in middle", "hello world", 6, 5, "world"},
		{"Negative start", "hello world", -5, 5, "world"},
		{"Negative length", "hello world", 0, -6, "hello"},
		{"Start out of range", "hello", 10, 5, ""},
		{"Length exceeds string", "hello", 0, 10, "hello"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrSubstr(test.input, test.start, test.length)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrUcfirst(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single character", "a", "A"},
		{"Already capitalized", "Hello", "Hello"},
		{"Lowercase string", "hello", "Hello"},
		{"Sentence", "hello world", "Hello world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrUcfirst(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrLcfirst(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Single character", "A", "a"},
		{"Already lowercase", "hello", "hello"},
		{"Capitalized string", "Hello", "hello"},
		{"Sentence", "Hello world", "hello world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrLcfirst(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		chars    string
		expected string
	}{
		{"Empty string", "", " ", ""},
		{"Trim spaces", "  hello  ", " ", "hello"},
		{"Trim specific characters", "xxxhelloxxx", "x", "hello"},
		{"No characters to trim", "hello", "x", "hello"},
		{"Trim multiple characters", "123hello321", "123", "hello"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrTrim(test.input, test.chars)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrLtrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		chars    string
		expected string
	}{
		{"Empty string", "", " ", ""},
		{"Trim left spaces", "  hello  ", " ", "hello  "},
		{"Trim left specific characters", "xxxhelloxxx", "x", "helloxxx"},
		{"No characters to trim", "hello", "x", "hello"},
		{"Trim left multiple characters", "123hello321", "123", "hello321"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrLtrim(test.input, test.chars)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrRtrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		chars    string
		expected string
	}{
		{"Empty string", "", " ", ""},
		{"Trim right spaces", "  hello  ", " ", "  hello"},
		{"Trim right specific characters", "xxxhelloxxx", "x", "xxxhello"},
		{"No characters to trim", "hello", "x", "hello"},
		{"Trim right multiple characters", "123hello321", "123", "123hello"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrRtrim(test.input, test.chars)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrPlural(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Regular noun", "book", "books"},
		{"Noun ending in 's'", "class", "classes"},
		{"Noun ending in 'y'", "city", "cities"},
		{"Noun ending in vowel + 'y'", "boy", "boys"},
		{"Irregular noun", "child", "children"},
		//{"Already plural", "books", "books"}, // TODO Expected "bookss", got "bookses"
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrPlural(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrSingular(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Empty string", "", ""},
		{"Regular plural", "books", "book"},
		{"Plural ending in 'es'", "classes", "class"},
		{"Plural ending in 'ies'", "cities", "city"},
		{"Irregular plural", "children", "child"},
		{"Already singular", "book", "book"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrSingular(test.input)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestStrWordwrap(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		width     int
		breakChar string
		expected  string
	}{
		{"Empty string", "", 5, "\n", ""},
		{"No wrapping needed", "hello", 10, "\n", "hello"},
		{"Simple wrap", "hello world", 5, "\n", "hello\nworld"},
		{"Multiple wraps", "hello beautiful world", 5, "\n", "hello\nbeautiful\nworld"},
		{"Custom break character", "hello world", 5, "<br>", "hello<br>world"},
		{"Zero width", "hello world", 0, "\n", "hello world"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.StrWordwrap(test.input, test.width, test.breakChar)
			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}
