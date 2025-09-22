package auth

import (
	"errors"
	"k2ray/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims defines the structure of the JWT claims for this application.
// It includes the standard RegisteredClaims and custom claims like Username.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateTokens creates both an access token and a refresh token for a given username.
func GenerateTokens(userID int64, username string) (accessToken string, refreshToken string, err error) {
	// Generate the access token (short-lived)
	accessToken, err = generateToken(userID, username, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	// Generate the refresh token (long-lived)
	refreshToken, err = generateToken(userID, username, 7*24*time.Hour) // 7 days
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken is a helper function to create a new JWT with a specific username and expiration.
func generateToken(userID int64, username string, expiration time.Duration) (string, error) {
	if config.AppConfig.JWTSecret == "" {
		return "", errors.New("JWT secret is not configured")
	}

	expirationTime := time.Now().Add(expiration)
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(), // JTI (JWT ID)
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "k2ray",
		},
	}

	// Create a new token object, specifying the signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	jwtKey := []byte(config.AppConfig.JWTSecret)
	return token.SignedString(jwtKey)
}

// ValidateToken parses a token string, validates its signature, and returns the claims.
func ValidateToken(tokenString string) (*Claims, error) {
	if config.AppConfig.JWTSecret == "" {
		return nil, errors.New("JWT secret is not configured")
	}

	claims := &Claims{}
	jwtKey := []byte(config.AppConfig.JWTSecret)

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
