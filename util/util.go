package util

import (
	"gounter/api/auth"
	"time"

	"github.com/golang-jwt/jwt"
)

// Helper function to generate a valid JWT token
func GenerateValidJWT() (string, error) {
	claims := &jwt.MapClaims{
		// this is the only thing we are validating
		// Set expiration to 5 minutes from now
		"exp": time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(auth.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
