package service

import (
	"fmt"
	"gfly/app/modules/auth"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"strings"
	"time"
)

// GenerateTokens func for generate a new Access & Refresh tokens.
func GenerateTokens(id string, credentials []string) (*auth.Tokens, error) {
	// Generate JWT Access token.
	accessToken, err := generateAccessToken(id, credentials)
	if err != nil {
		// Return token generation error.
		return nil, err
	}

	// Generate JWT Refresh token.
	refreshToken, err := generateRefreshToken()
	if err != nil {
		// Return refresh generation error.
		return nil, err
	}

	return &auth.Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateAccessToken(id string, credentials []string) (string, error) {
	// Get secret key from .env file.
	secret := utils.Getenv(auth.SecretKey, "")

	// Set expired minutes count for a secret key from .env file.
	ttlMinutes := utils.Getenv(auth.TtlMinutes, 0)

	// Create a new claims.
	claims := jwt.MapClaims{}

	// Set public claims:
	claims["id"] = id
	claims["expires"] = time.Now().Add(time.Minute * time.Duration(ttlMinutes)).Unix()

	// Set private token credentials:
	for _, credential := range credentials {
		claims[credential] = true
	}

	// Create a new JWT access token with claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate token.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// Return error, it JWT token generation failed.
		return "", err
	}

	return t, nil
}

func generateRefreshToken() (string, error) {
	hash := utils.Sha256(utils.Getenv(auth.RefreshKey, "") + time.Now().String())

	// Get expired days for refresh key from .env file.
	overDays := utils.Getenv(auth.TtlOverDays, 0)

	// Create expiration time.
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(overDays*24)).Unix())

	// Create a new refresh token (sha256 string with salt + expire time).
	t := hash + "." + expireTime

	return t, nil
}

// IsValidRefreshToken func for parse second argument from refresh token.
// A refresh token is valid is not expired.
func IsValidRefreshToken(refreshToken string) bool {
	tokenString := strings.Split(refreshToken, ".")
	if len(tokenString) < 2 {
		return false
	}
	expires, err := strconv.ParseInt(tokenString[1], 0, 64)
	if err != nil {
		log.Infof("parse refresh token error %v", err)

		return false
	}

	if expires < time.Now().Unix() {
		log.Info("refresh token expired")

		return false
	}

	return true
}
