package utils

import (
	cryptorand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// StrTruncateString truncates a string to the specified length and adds an ellipsis if truncated.
// It returns the original string if its length is less than or equal to maxLength,
// otherwise returns the truncated string with "..." appended.
//
// Parameters:
//   - s: The input string to truncate
//   - maxLength: The maximum allowed length of the string
//
// Returns:
//   - The truncated string with "..." appended if truncation occurred, otherwise original string
func StrTruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}

// StrSlugifyString converts a string to a URL-friendly slug.
// It performs the following transformations:
//   - Converts to lowercase
//   - Replaces spaces with hyphens
//   - Removes all special characters except letters, numbers and hyphens
//   - Replaces multiple hyphens with a single hyphen
//   - Trims hyphens from start and end
//
// Parameters:
//   - s: The input string to convert to slug
//
// Returns:
//   - A URL-friendly slug string
func StrSlugifyString(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	s = reg.ReplaceAllString(s, "")

	// Replace multiple hyphens with a single hyphen
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}

// StrGenerateRandomString generates a random string of the specified length.
//
// Parameters:
//   - length: The desired length of the random string
//
// Returns:
//   - string: The generated random string
//   - error: An error if random generation fails
func StrGenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := cryptorand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

// StrIsEmptyOrWhitespace checks if a string is empty or contains only whitespace characters.
//
// Parameters:
//   - s: The string to check
//
// Returns:
//   - bool: True if the string is empty or contains only whitespace, false otherwise
func StrIsEmptyOrWhitespace(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// StrContainsAny checks if a string contains any of the specified substrings.
//
// Parameters:
//   - s: The string to search in
//   - substrings: Variable number of substrings to search for
//
// Returns:
//   - bool: True if the string contains any of the substrings, false otherwise
func StrContainsAny(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// StrToTitleCase converts a string to title case format where the first letter
// of each word is capitalized.
//
// Parameters:
//   - s: The string to convert to title case
//
// Returns:
//   - string: The title cased string
func StrToTitleCase(s string) string {
	words := strings.Fields(strings.ToLower(s))
	for i, word := range words {
		if len(word) > 0 {
			r := []rune(word)
			r[0] = unicode.ToUpper(r[0])
			words[i] = string(r)
		}
	}
	return strings.Join(words, " ")
}

// StrRemoveNonAlphanumeric removes all non-alphanumeric characters from a string.
//
// Parameters:
//   - s: The string to process
//
// Returns:
//   - string: The string with only alphanumeric characters
func StrRemoveNonAlphanumeric(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]")
	return reg.ReplaceAllString(s, "")
}

// StrMaskString masks a portion of a string with the specified character.
//
// Parameters:
//   - s: The string to mask
//   - startVisible: Number of characters to leave visible at start
//   - endVisible: Number of characters to leave visible at end
//   - maskChar: The character to use for masking
//
// Returns:
//   - string: The masked string
//
// Example: StrMaskString("1234567890", 4, 4, '*') returns "1234****90"
func StrMaskString(s string, startVisible, endVisible int, maskChar rune) string {
	if len(s) <= startVisible+endVisible {
		return s
	}

	start := s[:startVisible]
	end := s[len(s)-endVisible:]
	masked := strings.Repeat(string(maskChar), len(s)-startVisible-endVisible)

	return start + masked + end
}

// StrPadLeft pads a string on the left side with a specified character to reach
// the desired length.
//
// Parameters:
//   - s: The string to pad
//   - padChar: The character to use for padding
//   - length: The desired total length
//
// Returns:
//   - string: The padded string
func StrPadLeft(s string, padChar rune, length int) string {
	if len(s) >= length {
		return s
	}

	padding := strings.Repeat(string(padChar), length-len(s))
	return padding + s
}

// StrPadRight pads a string on the right side with a specified character to reach
// the desired length.
//
// Parameters:
//   - s: The string to pad
//   - padChar: The character to use for padding
//   - length: The desired total length
//
// Returns:
//   - string: The padded string
func StrPadRight(s string, padChar rune, length int) string {
	if len(s) >= length {
		return s
	}

	padding := strings.Repeat(string(padChar), length-len(s))
	return s + padding
}

// StrReverseString reverses the characters in a string.
//
// Parameters:
//   - s: The string to reverse
//
// Returns:
//   - string: The reversed string
func StrReverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// StrCountWords counts the number of words in a string.
// Words are considered to be separated by whitespace.
//
// Parameters:
//   - s: The string to count words in
//
// Returns:
//   - int: The number of words in the string
func StrCountWords(s string) int {
	if StrIsEmptyOrWhitespace(s) {
		return 0
	}

	words := strings.Fields(s)
	return len(words)
}

// StrTruncateWords truncates a string to the specified number of words and adds
// an ellipsis if the string was truncated.
//
// Parameters:
//   - s: The string to truncate
//   - maxWords: Maximum number of words to keep
//
// Returns:
//   - string: The truncated string
func StrTruncateWords(s string, maxWords int) string {
	if maxWords <= 0 {
		return ""
	}

	words := strings.Fields(s)
	if len(words) <= maxWords {
		return s
	}

	return strings.Join(words[:maxWords], " ") + "..."
}

// StrFormatWithCommas formats a number as a string with commas as thousand separators.
//
// Parameters:
//   - n: The number to format
//
// Returns:
//   - string: The formatted number string
func StrFormatWithCommas(n int64) string {
	return fmt.Sprintf("%d", n)
}

// StrCountRunes counts the number of Unicode characters (runes) in a string.
//
// Parameters:
//   - s: The string to count runes in
//
// Returns:
//   - int: The number of runes in the string
func StrCountRunes(s string) int {
	return utf8.RuneCountInString(s)
}

// StrCamelCase converts a string to camelCase format where the first word is lowercase
// and subsequent words are capitalized with no spaces.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - string: The camelCase formatted string
func StrCamelCase(s string) string {
	words := strings.Fields(strings.ToLower(s))
	if len(words) == 0 {
		return ""
	}

	result := words[0]
	for i := 1; i < len(words); i++ {
		if len(words[i]) > 0 {
			result += strings.Title(words[i])
		}
	}

	return result
}

// StrSnakeCase converts a string to snake_case format.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - string: The snake_case formatted string
func StrSnakeCase(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with underscores
	s = strings.ReplaceAll(s, " ", "_")

	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9_]")
	s = reg.ReplaceAllString(s, "")

	// Replace multiple underscores with a single underscore
	reg = regexp.MustCompile("_+")
	s = reg.ReplaceAllString(s, "_")

	// Trim underscores from start and end
	s = strings.Trim(s, "_")

	return s
}

// StrKebabCase converts a string to kebab-case format.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - string: The kebab-case formatted string
func StrKebabCase(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces with hyphens
	s = strings.ReplaceAll(s, " ", "-")

	// Remove special characters
	reg := regexp.MustCompile("[^a-z0-9-]")
	s = reg.ReplaceAllString(s, "")

	// Replace multiple hyphens with a single hyphen
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from start and end
	s = strings.Trim(s, "-")

	return s
}

// StrAfter returns the portion of a string after the first occurrence of a given value.
//
// Parameters:
//   - s: The string to search in
//   - search: The substring to search for
//
// Returns:
//   - string: Everything after the search string, or the entire string if not found
func StrAfter(s, search string) string {
	if search == "" {
		return s
	}

	pos := strings.Index(s, search)
	if pos == -1 {
		return s
	}

	return s[pos+len(search):]
}

// StrAfterLast returns the portion of a string after the last occurrence of a given value.
//
// Parameters:
//   - s: The string to search in
//   - search: The substring to search for
//
// Returns:
//   - string: Everything after the last occurrence of search string, or entire string if not found
func StrAfterLast(s, search string) string {
	if search == "" {
		return s
	}

	pos := strings.LastIndex(s, search)
	if pos == -1 {
		return s
	}

	return s[pos+len(search):]
}

// StrBefore returns the portion of a string before the first occurrence of a given value.
//
// Parameters:
//   - s: The string to search in
//   - search: The substring to search for
//
// Returns:
//   - string: Everything before the search string, or the entire string if not found
func StrBefore(s, search string) string {
	if search == "" {
		return s
	}

	pos := strings.Index(s, search)
	if pos == -1 {
		return s
	}

	return s[:pos]
}

// StrBeforeLast returns the portion of a string before the last occurrence of a given value.
//
// Parameters:
//   - s: The string to search in
//   - search: The substring to search for
//
// Returns:
//   - string: Everything before the last occurrence of search string, or entire string if not found
func StrBeforeLast(s, search string) string {
	if search == "" {
		return s
	}

	pos := strings.LastIndex(s, search)
	if pos == -1 {
		return s
	}

	return s[:pos]
}

// StrBetween returns the portion of a string between two values.
//
// Parameters:
//   - s: The string to search in
//   - start: The starting substring
//   - end: The ending substring
//
// Returns:
//   - string: The portion between start and end strings, or entire string if not found
func StrBetween(s, start, end string) string {
	if start == "" || end == "" {
		return s
	}

	startPos := strings.Index(s, start)
	if startPos == -1 {
		return s
	}

	subStr := s[startPos+len(start):]
	endPos := strings.Index(subStr, end)
	if endPos == -1 {
		return s
	}

	return subStr[:endPos]
}

// StrContains determines if a string contains a given substring.
//
// Parameters:
//   - s: The string to search in
//   - substr: The substring to search for
//
// Returns:
//   - bool: True if substring is found, false otherwise
func StrContains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// StrContainsAll determines if a string contains all of the given substrings.
//
// Parameters:
//   - s: The string to search in
//   - substrings: Variable number of substrings to search for
//
// Returns:
//   - bool: True if all substrings are found, false otherwise
func StrContainsAll(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}

// StrEndsWith determines if a string ends with a given substring.
//
// Parameters:
//   - s: The string to check
//   - substrings: Variable number of possible ending substrings
//
// Returns:
//   - bool: True if string ends with any of the substrings, false otherwise
func StrEndsWith(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.HasSuffix(s, substr) {
			return true
		}
	}
	return false
}

// StrFinish appends a single instance of the given value to a string
// if it does not already end with it.
//
// Parameters:
//   - s: The string to append to
//   - cap: The string to append
//
// Returns:
//   - string: The resulting string
func StrFinish(s, cap string) string {
	if cap == "" {
		return s
	}

	if strings.HasSuffix(s, cap) {
		return s
	}

	return s + cap
}

// StrIs determines if a string matches a given pattern.
// Asterisks may be used as wildcard values.
//
// Parameters:
//   - pattern: The pattern to match against (can include * wildcards)
//   - s: The string to check
//
// Returns:
//   - bool: True if string matches pattern, false otherwise
func StrIs(pattern, s string) bool {
	if pattern == s {
		return true
	}

	// Convert the pattern to a regular expression
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = "^" + pattern + "$"

	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

// StrIsAscii determines if a string contains only 7-bit ASCII characters.
//
// Parameters:
//   - s: The string to check
//
// Returns:
//   - bool: True if string contains only ASCII characters, false otherwise
func StrIsAscii(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// StrLength returns the number of Unicode characters in a string.
//
// Parameters:
//   - s: The string to measure
//
// Returns:
//   - int: The length in characters (runes)
func StrLength(s string) int {
	return utf8.RuneCountInString(s)
}

// StrLimit truncates a string to the specified length.
//
// Parameters:
//   - s: The string to truncate
//   - limit: Maximum length
//
// Returns:
//   - string: The truncated string
func StrLimit(s string, limit int) string {
	return StrTruncateString(s, limit)
}

// StrLower converts a string to lowercase.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - string: The lowercase string
func StrLower(s string) string {
	return strings.ToLower(s)
}

// StrUpper converts a string to uppercase.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - string: The uppercase string
func StrUpper(s string) string {
	return strings.ToUpper(s)
}

// StrRandom generates a random string of specified length.
//
// Parameters:
//   - length: The desired length of the random string
//
// Returns:
//   - string: The generated random string
func StrRandom(length int) string {
	// Initialize the random number generator with a seed
	rand.Seed(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// StrReplace replaces all occurrences of a given value in a string with another value.
//
// Parameters:
//   - search: The string to find
//   - replace: The string to replace with
//   - subject: The string to perform replacements on
//
// Returns:
//   - string: The resulting string after replacements
func StrReplace(search, replace, subject string) string {
	return strings.ReplaceAll(subject, search, replace)
}

// StrReplaceArray replaces a search string with an array of replacements sequentially.
//
// Parameters:
//   - search: The string to find
//   - replace: Array of replacement strings
//   - subject: The string to perform replacements on
//
// Returns:
//   - string: The resulting string after replacements
func StrReplaceArray(search string, replace []string, subject string) string {
	result := subject
	for _, value := range replace {
		pos := strings.Index(result, search)
		if pos == -1 {
			break
		}
		result = result[:pos] + value + result[pos+len(search):]
	}
	return result
}

// StrReplaceFirst replaces the first occurrence of a given value in a string.
//
// Parameters:
//   - search: The string to find
//   - replace: The string to replace with
//   - subject: The string to perform replacement on
//
// Returns:
//   - string: The resulting string after replacement
func StrReplaceFirst(search, replace, subject string) string {
	pos := strings.Index(subject, search)
	if pos == -1 {
		return subject
	}

	return subject[:pos] + replace + subject[pos+len(search):]
}

// StrReplaceLast replaces the last occurrence of a given value in a string.
//
// Parameters:
//   - search: The string to find
//   - replace: The string to replace with
//   - subject: The string to perform replacement on
//
// Returns:
//   - string: The resulting string after replacement
func StrReplaceLast(search, replace, subject string) string {
	pos := strings.LastIndex(subject, search)
	if pos == -1 {
		return subject
	}

	return subject[:pos] + replace + subject[pos+len(search):]
}

// StrStart prepends a value to a string if it doesn't already start with it.
//
// Parameters:
//   - s: The string to prepend to
//   - prefix: The string to prepend
//
// Returns:
//   - string: The resulting string
func StrStart(s, prefix string) string {
	if prefix == "" {
		return s
	}

	if strings.HasPrefix(s, prefix) {
		return s
	}

	return prefix + s
}

// StrStartsWith determines if a string starts with any of the given substrings.
//
// Parameters:
//   - s: The string to check
//   - substrings: Variable number of possible starting substrings
//
// Returns:
//   - bool: True if string starts with any substring, false otherwise
func StrStartsWith(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if strings.HasPrefix(s, substr) {
			return true
		}
	}
	return false
}

// StrStudly converts a string to StudlyCase format.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - string: The StudlyCase formatted string
func StrStudly(s string) string {
	// Replace hyphens and underscores with spaces
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	// Convert to title case
	s = StrToTitleCase(s)

	// Remove spaces
	return strings.ReplaceAll(s, " ", "")
}

// StrSubstr returns a portion of a string based on start position and length.
//
// Parameters:
//   - s: The string to get a substring from
//   - start: Starting position
//   - length: Length of substring
//
// Returns:
//   - string: The substring
func StrSubstr(s string, start, length int) string {
	runes := []rune(s)
	l := len(runes)

	// Handle negative start
	if start < 0 {
		start = l + start
		if start < 0 {
			start = 0
		}
	}

	// Handle out of range start
	if start >= l {
		return ""
	}

	// Handle negative length
	if length < 0 {
		length = l - start + length
		if length < 0 {
			length = 0
		}
	}

	// Handle out of range length
	if start+length > l {
		length = l - start
	}

	return string(runes[start : start+length])
}

// StrUcfirst capitalizes the first character of a string.
//
// Parameters:
//   - s: The string to capitalize
//
// Returns:
//   - string: The string with first character capitalized
func StrUcfirst(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// StrLcfirst converts the first character of a string to lowercase.
//
// Parameters:
//   - s: The string to convert
//
// Returns:
//   - string: The string with first character lowercased
func StrLcfirst(s string) string {
	if s == "" {
		return ""
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// StrTrim removes specified characters from both ends of a string.
//
// Parameters:
//   - s: The string to trim
//   - chars: The characters to remove
//
// Returns:
//   - string: The trimmed string
func StrTrim(s, chars string) string {
	return strings.Trim(s, chars)
}

// StrLtrim removes specified characters from the start of a string.
//
// Parameters:
//   - s: The string to trim
//   - chars: The characters to remove
//
// Returns:
//   - string: The left-trimmed string
func StrLtrim(s, chars string) string {
	return strings.TrimLeft(s, chars)
}

// StrRtrim removes specified characters from the end of a string.
//
// Parameters:
//   - s: The string to trim
//   - chars: The characters to remove
//
// Returns:
//   - string: The right-trimmed string
func StrRtrim(s, chars string) string {
	return strings.TrimRight(s, chars)
}

// StrPlural converts a singular word to its plural form.
// This is a simple implementation and may not work for all cases.
//
// Parameters:
//   - s: The singular word to pluralize
//
// Returns:
//   - string: The plural form of the word
func StrPlural(s string) string {
	if s == "" {
		return ""
	}

	// Some common irregular plurals
	irregulars := map[string]string{
		"child":  "children",
		"goose":  "geese",
		"man":    "men",
		"woman":  "women",
		"tooth":  "teeth",
		"foot":   "feet",
		"mouse":  "mice",
		"person": "people",
	}

	if plural, ok := irregulars[strings.ToLower(s)]; ok {
		return plural
	}

	// Handle words ending in 'y'
	if strings.HasSuffix(s, "y") {
		// If the word ends in a vowel + y, just add 's'
		if len(s) > 1 && strings.Contains("aeiou", string(s[len(s)-2])) {
			return s + "s"
		}
		// Otherwise, replace 'y' with 'ies'
		return s[:len(s)-1] + "ies"
	}

	// Handle words ending in 's', 'x', 'z', 'ch', 'sh'
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "x") || strings.HasSuffix(s, "z") ||
		strings.HasSuffix(s, "ch") || strings.HasSuffix(s, "sh") {
		return s + "es"
	}

	// Default: just add 's'
	return s + "s"
}

// StrSingular converts a plural word to its singular form.
// This is a simple implementation and may not work for all cases.
//
// Parameters:
//   - s: The plural word to singularize
//
// Returns:
//   - string: The singular form of the word
func StrSingular(s string) string {
	if s == "" {
		return ""
	}

	// Some common irregular singulars
	irregulars := map[string]string{
		"children": "child",
		"geese":    "goose",
		"men":      "man",
		"women":    "woman",
		"teeth":    "tooth",
		"feet":     "foot",
		"mice":     "mouse",
		"people":   "person",
	}

	if singular, ok := irregulars[strings.ToLower(s)]; ok {
		return singular
	}

	// Handle words ending in 'ies'
	if strings.HasSuffix(s, "ies") {
		return s[:len(s)-3] + "y"
	}

	// Handle words ending in 'es'
	if strings.HasSuffix(s, "es") {
		// Check if it's one of the special cases
		base := s[:len(s)-2]
		if strings.HasSuffix(base, "s") || strings.HasSuffix(base, "x") || strings.HasSuffix(base, "z") ||
			strings.HasSuffix(base, "ch") || strings.HasSuffix(base, "sh") {
			return base
		}
	}

	// Handle words ending in 's'
	if strings.HasSuffix(s, "s") {
		return s[:len(s)-1]
	}

	// Default: return as is
	return s
}

// StrWordwrap wraps a string to a given number of characters.
//
// Parameters:
//   - s: The string to wrap
//   - width: The number of characters at which to wrap
//   - breakChar: The string to insert at break points
//
// Returns:
//   - string: The wrapped string
func StrWordwrap(s string, width int, breakChar string) string {
	if width <= 0 {
		return s
	}

	var result strings.Builder
	words := strings.Fields(s)

	lineLength := 0
	for i, word := range words {
		wordLength := len(word)

		if i > 0 {
			if lineLength+wordLength+1 > width {
				result.WriteString(breakChar)
				lineLength = 0
			} else {
				result.WriteString(" ")
				lineLength++
			}
		}

		result.WriteString(word)
		lineLength += wordLength
	}

	return result.String()
}
