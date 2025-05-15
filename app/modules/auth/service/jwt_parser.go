package service

import (
	"gfly/app/modules/auth"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"strings"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      int
	Credentials core.Data
	Expires     int64
}

// ExtractToken func to get JWT from header.
func ExtractToken(c *core.Ctx) string {
	bearToken := c.Root().Request.Header.Peek(core.HeaderAuthorization)

	// Normally Authorization HTTP header.
	onlyToken := strings.Split(utils.UnsafeStr(bearToken), " ")
	if len(onlyToken) == 2 {
		return onlyToken[1]
	}

	return ""
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(tokenString string) (*TokenMetadata, error) {
	token, err := verifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userID, _ := strconv.Atoi(claims["id"].(string))

		expires := int64(claims["expires"].(float64))

		credentials := make(core.Data)

		return &TokenMetadata{
			UserID:      userID,
			Credentials: credentials,
			Expires:     expires,
		}, nil
	}
	return nil, err
}

// verifyToken function will parse, validate and verify the signature
func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// jwtKeyFunc will receive the parsed token and should return the cryptographic key
// for verifying the signature
func jwtKeyFunc(_ *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv(auth.SecretKey)), nil
}
