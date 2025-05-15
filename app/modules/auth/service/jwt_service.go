package service

import (
	"fmt"
	"gfly/app/domain/models"
	"gfly/app/domain/models/types"
	"gfly/app/domain/repository"
	"gfly/app/modules/auth"
	"gfly/app/modules/auth/dto"
	"github.com/gflydev/cache"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	mb "github.com/gflydev/db"
	"github.com/gflydev/db/null"
	"strconv"
	"strings"
	"time"
)

// SignIn login app.
func SignIn(signIn *dto.SignIn) (*auth.Tokens, error) {
	// Get user by email.
	user := repository.Pool.GetUserByEmail(signIn.Username)
	if user == nil {
		return nil, errors.New("Invalid email address or password")
	}
	// Compare given user password with stored in found user.
	isValidPassword := utils.ComparePasswords(user.Password, signIn.Password)
	if !isValidPassword {
		return nil, errors.New("Invalid email address or password")
	}

	userIDStr := strconv.Itoa(user.ID)
	// Generate a new pair of access and refresh tokens.
	tokens, err := GenerateTokens(userIDStr, make([]string, 0))
	if err != nil {
		log.Errorf("Error while generating tokens %q", err)
		return nil, err
	}

	// Set expired days from .env file
	ttlDays := utils.Getenv(auth.TtlOverDays, 0) // 7 days by default

	// Save refresh token to Redis.
	expiredTime := time.Duration(ttlDays*24*3600) * time.Second // 604 800 seconds = 7 days
	if err = cache.Set(userIDStr, tokens.Refresh, expiredTime); err != nil {
		log.Errorf("Error while caching to token to Redis %q", err)
		return nil, err
	}

	return tokens, nil
}

// SignUp register new user.
func SignUp(signUp *dto.SignUp) (*models.User, error) {
	email := strings.ToLower(signUp.Email)

	userEmail := repository.Pool.GetUserByEmail(email)
	if userEmail != nil {
		return nil, errors.New("User with the given email address already exists")
	}

	// Create a new user struct.
	user := &models.User{}

	// Set initialized default data for user
	user.Email = email
	user.Password = utils.GeneratePassword(signUp.Password)
	user.Fullname = signUp.Fullname
	user.Phone = signUp.Phone
	user.Token = null.String("")
	user.Status = types.UserStatusActive
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.LastAccessAt = null.NowTime()

	// Create a new user with validated data.
	err := mb.CreateModel(user)
	if err != nil {
		log.Errorf("Error while creating new user %q with data '%v'", err, user)
		return nil, errors.New("Error occurs while signup user")
	}

	return user, nil
}

// SignOut function takes in jwtToken string, utils.ExtractTokenMetadata extract access token metadata
// to get a userID which is the key that store refresh token in the Redis Caching
// then delete refresh token from the Redis
// and DeleteToken will delete access token
// by send it to black-list (middleware will handle invalid token in blacklist).
func SignOut(jwtToken string) error {
	// Extract access token metadata
	claims, err := ExtractTokenMetadata(jwtToken)
	if err != nil {
		log.Errorf("Error while logging out %q", err)
		return errors.New("Logout error")
	}

	userID := strconv.Itoa(claims.UserID)

	// Delete refresh token from Redis.
	if err = cache.Del(userID); err != nil {
		log.Errorf("Error while delete refresh token from Redis %q", err)
		return err
	}

	// Delete access token by send it to black-list
	DeleteToken(jwtToken)
	return nil
}

// RefreshToken function to refresh JWT token from user.
func RefreshToken(jwtToken, refreshToken string) (*auth.Tokens, error) {
	// Get claims from JWT.
	claims, err := ExtractTokenMetadata(jwtToken)
	if err != nil {
		log.Errorf("Error while extracting token metadata %q", err)
		return nil, errors.New("Refresh token error")
	}
	// Define user ID.
	userID := claims.UserID
	userIDStr := strconv.Itoa(userID)

	// Get refresh token from Redis.
	val, err := cache.Get(userIDStr)
	if err != nil {
		log.Errorf("Error while getting refresh token from Redis %q", err)
		return nil, errors.New("Refresh token error")
	}

	if refreshToken != val {
		return nil, errors.New("Refresh token mismatch")
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := GenerateTokens(userIDStr, make([]string, 0))
	if err != nil {
		log.Errorf("Error while generating JWT Tokens")
		return nil, errors.New("Refresh token error")
	}

	// Set expired days from .env file.
	ttlDays := utils.Getenv(auth.TtlOverDays, 0)
	duration := time.Duration(ttlDays*7*24*3600) * time.Second

	// Update refresh token to Redis.
	if err = cache.Set(userIDStr, tokens.Refresh, duration); err != nil {
		log.Errorf("Refresh token error '%v'", err)

		return nil, errors.New("Refresh token error")
	}

	// Delete JWT token by sending it to blacklist
	DeleteToken(jwtToken)

	return tokens, nil
}

// DeleteToken add jwtToken to blacklist
func DeleteToken(jwtToken string) bool {
	key := fmt.Sprintf("%s:%s", utils.Getenv(auth.Blacklist, ""), jwtToken)

	// Set expired minutes count for a secret key from .env file.
	ttlMinutes := utils.Getenv(auth.TtlMinutes, 0)
	expiresTime := time.Duration(ttlMinutes*60) * time.Second

	// Update refresh token to Redis.
	if err := cache.Set(key, "blocked", expiresTime); err != nil {
		log.Errorf("Delete JWT token error '%v'", err)

		return false
	}

	return true
}

// IsBlockedToken Check if jwtToken is locked or not
func IsBlockedToken(jwtToken string) (bool, error) {
	isCheckBlacklist := utils.Getenv(auth.CheckBlacklist, false)
	if !isCheckBlacklist {
		return false, nil
	}

	key := fmt.Sprintf("%s:%s", utils.Getenv(auth.Blacklist, ""), jwtToken)

	// Get blocked JWT in Redis.
	val, err := cache.Get(key)
	if err != nil {
		return false, nil
	}
	exists := val == string(types.UserStatusBlocked)

	return exists, nil
}
