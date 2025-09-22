package services

import (
	"fmt"
	"gfly/internal/domain/models"
	"gfly/internal/domain/models/types"
	"gfly/internal/domain/repository"
	"gfly/internal/modules/auth"
	"gfly/internal/modules/auth/dto"
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

// ====================================================================
// ========================= Main functions ===========================
// ====================================================================

// SignIn authenticates a user and generates access/refresh token pair
//
// Parameters:
//   - signIn: *dto.SignIn - Contains validated login credentials:
//   - Username: Email address used for login
//   - Password: Plain text password to validate
//
// Returns:
//   - *auth.Token: Token pair containing access and refresh tokens if successful
//   - error: Error if authentication fails:
//   - Invalid email/password
//   - User account not active
//   - Token generation failed
//   - Redis caching failed
//
// Flow:
// 1. Looks up user by email address
// 2. Validates provided password against stored hash
// 3. Verifies user account is active
// 4. Generates new access/refresh token pair
// 5. Caches refresh token in Redis with TTL
//
// Example:
//
//	 credentials := &dto.SignIn{
//		Username: "user@example.com",
//		Password: "secret123"
//	 }
//	 tokens, err := SignIn(credentials)
func SignIn(signIn dto.SignIn) (*auth.Token, error) {
	// Get user by email.
	user := repository.Pool.GetUserByEmail(signIn.Username)
	if user == nil {
		return nil, errors.New("Invalid email address or password")
	}
	// Compare a given user password with stored in found user.
	isValidPassword := utils.ComparePasswords(user.Password, signIn.Password)
	if !isValidPassword {
		return nil, errors.New("Invalid email address or password")
	}

	if user.Status != types.UserStatusActive {
		return nil, errors.New("User is not activated")
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

	// Save the refresh token to Redis.
	expiredTime := time.Duration(ttlDays*24*3600) * time.Second // 604 800 seconds = 7 days
	if err = cache.Set(userIDStr, tokens.Refresh, expiredTime); err != nil {
		log.Errorf("Error while caching to token to Redis %q", err)
		return nil, err
	}

	return tokens, nil
}

// SignUp creates a new user account with the provided signup details.
//
// Parameters:
//   - signUp: *dto.SignUp - Contains validated user registration data including:
//   - Email: User's email address (will be converted to lowercase)
//   - Password: Plain text password that will be hashed
//   - Fullname: User's full name
//   - Phone: User's phone number
//   - Avatar: Optional profile image URL
//   - Status: Optional account status
//
// Returns:
//   - *models.User: Newly created user record if successful
//   - error: Error if user creation fails:
//   - Email already exists
//   - Database errors during user creation
//
// Flow:
// 1. Converts email to lowercase
// 2. Checks if email is already registered
// 3. Creates new user with provided details:
//   - Hashes the password
//   - Sets default status to active
//   - Sets creation/update timestamps
//
// 4. Saves user to database
//
// Example:
//
//	 signup := &dto.SignUp{
//		Email: "user@example.com",
//		Password: "secret123",
//		Fullname: "John Doe",
//		Phone: "1234567890"
//	 }
//	 user, err := SignUp(signup)
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
	user.LastAccessAt = null.TimeNow()

	// Create a new user with validated data.
	err := mb.CreateModel(user)
	if err != nil {
		log.Errorf("Error while creating new user %q with data '%v'", err, user)
		return nil, errors.New("Error occurs while signup user")
	}

	return user, nil
}

// SignOut handles user logout by invalidating both refresh and access tokens
//
// Parameters:
//   - jwtToken: The current access token to be invalidated
//
// Returns:
//   - error: Error if token invalidation fails
//
// Flow:
// 1. Extracts user ID and metadata from the access token
// 2. Uses the user ID to find and delete the refresh token from Redis
// 3. Adds the access token to the blacklist to invalidate it
//
// Note that this implements a "logout everywhere" approach by:
// - Deleting the refresh token to prevent getting new access tokens
// - Blacklisting the current access token to immediately invalidate it
//
// Errors:
//   - Returns error if access token metadata extraction fails
//   - Returns error if refresh token deletion from Redis fails
//   - Continues execution if blacklisting access token fails (best effort)
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
	deleteToken(jwtToken)
	return nil
}

// RefreshToken creates new tokens by validating the existing access and refresh tokens.
//
// Parameters:
//   - jwtToken: The current access token to be refreshed
//   - refreshToken: The current refresh token to validate against stored token
//
// Returns:
//   - *auth.Token: New token pair containing fresh access and refresh tokens
//   - error: Error if token validation fails or token generation encounters issues
//
// Flow:
// 1. Extracts user ID and metadata from the access token
// 2. Validates the provided refresh token matches the one stored in Redis for the user
// 3. Generates new access and refresh token pair
// 4. Updates the new refresh token in Redis with TTL
// 5. Blacklists the old access token
//
// Errors:
//   - Returns error if access token metadata extraction fails
//   - Returns error if refresh token validation against Redis fails
//   - Returns error if refresh tokens don't match
//   - Returns error if generating new tokens fails
//   - Returns error if storing new refresh token in Redis fails
func RefreshToken(jwtToken, refreshToken string) (*auth.Token, error) {
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
		log.Errorf("Error while generating JWT Token")
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
	deleteToken(jwtToken)

	return tokens, nil
}

// IsBlockedToken checks if a JWT token has been blacklisted/blocked
//
// Parameters:
//   - jwtToken: The JWT token string to check
//
// Returns:
//   - bool: true if token is blocked, false otherwise
//   - error: Error if any issues occurred during check
//
// Flow:
// 1. Checks if blacklist checking is enabled in config
// 2. If disabled, returns false immediately
// 3. Constructs Redis key by combining blacklist prefix with JWT token
// 4. Queries Redis to check if token exists in blacklist
// 5. Returns true if token value matches blocked status
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

// ====================================================================
// ======================== Helper Functions ==========================
// ====================================================================

// deleteToken adds the JWT token to a blacklist in Redis cache to invalidate it
//
// Parameters:
//   - jwtToken: The JWT access token to be blacklisted
//
// Returns:
//   - bool: true if token was successfully blacklisted, false if there was an error
//
// Flow:
// 1. Constructs Redis key by combining blacklist prefix with JWT token
// 2. Gets TTL duration from environment config
// 3. Adds token to Redis blacklist with TTL
//
// The blacklisted token will automatically expire after the configured TTL period.
// This prevents the accumulation of expired blacklist entries in Redis.
//
// Errors:
//   - Returns false if setting token in Redis cache fails
//   - Logs error details but continues execution
func deleteToken(jwtToken string) bool {
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
