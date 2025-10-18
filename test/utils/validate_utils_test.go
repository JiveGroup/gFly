package utils

import (
	"gfly/pkg/utils"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"ValidEmail", "test@example.com", true},
		{"ValidEmailWithSubdomain", "test@sub.example.com", true},
		{"ValidEmailWithPlus", "test+tag@example.com", true},
		{"InvalidEmailNoAt", "testexample.com", false},
		{"InvalidEmailNoDomain", "test@", false},
		{"InvalidEmailNoUsername", "@example.com", false},
		{"InvalidEmailSpaces", "test @example.com", false},
		{"EmptyEmail", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.IsValidEmail(test.email)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for email: %s", test.expected, result, test.email)
			}
		})
	}
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"ValidHTTP", "http://example.com", true},
		{"ValidHTTPS", "https://example.com", true},
		{"ValidWithPath", "https://example.com/path", true},
		{"ValidWithQuery", "https://example.com/path?query=value", true},
		{"ValidWithPort", "https://example.com:8080", true},
		{"InvalidNoScheme", "example.com", false},
		{"InvalidNoHost", "http://", false},
		{"InvalidScheme", "invalid://example.com", true}, // Note: url.Parse accepts any scheme
		{"EmptyURL", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.IsValidURL(test.url)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for URL: %s", test.expected, result, test.url)
			}
		})
	}
}

func TestIsValidPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected bool
	}{
		{"ValidSimple", "1234567890", true},
		{"ValidWithDashes", "123-456-7890", true},
		{"ValidWithSpaces", "123 456 7890", true},
		{"ValidWithParentheses", "(123) 456-7890", true},
		{"ValidWithPlus", "+11234567890", true},
		{"ValidInternational", "+44 20 7946 0958", true},
		{"InvalidTooShort", "123456", false},
		{"InvalidTooLong", "12345678901234567890", false},
		{"InvalidWithLetters", "123-456-789a", false},
		{"EmptyPhone", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.IsValidPhoneNumber(test.phone)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for phone: %s", test.expected, result, test.phone)
			}
		})
	}
}

func TestIsValidIPv4(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"ValidSimple", "192.168.1.1", true},
		{"ValidZeros", "0.0.0.0", true},
		{"ValidMax", "255.255.255.255", true},
		{"InvalidTooFewParts", "192.168.1", false},
		{"InvalidTooManyParts", "192.168.1.1.1", false},
		{"InvalidNegative", "-1.168.1.1", false},
		{"InvalidTooLarge", "256.168.1.1", false},
		{"InvalidNonNumeric", "a.168.1.1", false},
		{"EmptyIP", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.IsValidIPv4(test.ip)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for IP: %s", test.expected, result, test.ip)
			}
		})
	}
}

func TestIsValidCreditCard(t *testing.T) {
	tests := []struct {
		name     string
		card     string
		expected bool
	}{
		{"ValidVisa", "4111111111111111", true},
		{"ValidMasterCard", "5555555555554444", true},
		{"ValidAmex", "378282246310005", true},
		{"ValidDiscover", "6011111111111117", true},
		{"ValidWithSpaces", "4111 1111 1111 1111", true},
		{"ValidWithDashes", "4111-1111-1111-1111", true},
		{"InvalidChecksum", "4111111111111112", false},
		{"InvalidTooShort", "41111111", false},
		{"InvalidTooLong", "41111111111111111111", false},
		{"InvalidNonNumeric", "411111111111111a", false},
		{"EmptyCard", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.IsValidCreditCard(test.card)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for card: %s", test.expected, result, test.card)
			}
		})
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"ValidStrong", "Abcd1234!", true},
		{"ValidComplex", "P@ssw0rd", true},
		{"InvalidTooShort", "Ab1!", false},
		{"InvalidNoUppercase", "abcd1234!", false},
		{"InvalidNoLowercase", "ABCD1234!", false},
		{"InvalidNoDigit", "Abcdefgh!", false},
		{"InvalidNoSpecial", "Abcd1234", false},
		{"EmptyPassword", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.IsStrongPassword(test.password)
			if result != test.expected {
				t.Errorf("Expected %v, got %v for password: %s", test.expected, result, test.password)
			}
		})
	}
}
