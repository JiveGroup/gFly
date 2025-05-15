package utils

import (
	"gfly/app/http/response"
	"github.com/gflydev/core"
	"github.com/gflydev/validation"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// Validate perform data input checking.
func Validate(structData any, msgForTagFunc ...validation.MsgForTagFunc) *response.Error {
	errorData, err := validation.Check(structData, msgForTagFunc...)

	if err != nil {
		// Convert validation data to core.Data
		data := make(core.Data)
		for attr := range errorData {
			data[attr] = errorData[attr]
		}

		// Response validation error
		return &response.Error{
			Code:    core.StatusBadRequest,
			Message: "Invalid input",
			Data:    data,
		}
	}

	return nil
}

// IsValidEmail checks if a string is a valid email address
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// IsValidURL checks if a string is a valid URL
func IsValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsValidPhoneNumber checks if a string is a valid phone number
// This is a simple implementation that checks for digits and common separators
func IsValidPhoneNumber(phone string) bool {
	// Remove common separators
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// Check if the result contains only digits
	match, _ := regexp.MatchString("^[0-9]+$", phone)
	return match && len(phone) >= 7 && len(phone) <= 15
}

// IsValidIPv4 checks if a string is a valid IPv4 address
func IsValidIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// IsValidCreditCard checks if a string is a valid credit card number using the Luhn algorithm
func IsValidCreditCard(cardNumber string) bool {
	// Remove spaces and dashes
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	// Check if the card number contains only digits
	match, _ := regexp.MatchString("^[0-9]+$", cardNumber)
	if !match {
		return false
	}

	// Check length (most credit cards are between 13 and 19 digits)
	if len(cardNumber) < 13 || len(cardNumber) > 19 {
		return false
	}

	// Luhn algorithm
	sum := 0
	isSecondDigit := false
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit := int(cardNumber[i] - '0')

		if isSecondDigit {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		isSecondDigit = !isSecondDigit
	}

	return sum%10 == 0
}

// IsStrongPassword checks if a password meets strong password criteria
func IsStrongPassword(password string) bool {
	// Check length
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter
	hasUpper, _ := regexp.MatchString("[A-Z]", password)
	if !hasUpper {
		return false
	}

	// Check for at least one lowercase letter
	hasLower, _ := regexp.MatchString("[a-z]", password)
	if !hasLower {
		return false
	}

	// Check for at least one digit
	hasDigit, _ := regexp.MatchString("[0-9]", password)
	if !hasDigit {
		return false
	}

	// Check for at least one special character
	hasSpecial, _ := regexp.MatchString("[^A-Za-z0-9]", password)

	return hasSpecial
}
