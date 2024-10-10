package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Define a secret key for signing the JWT (in a real application, store it securely)
var JWTSecret = []byte("my-secret-key")

const bearerPrefix = "Bearer "

// AuthorizationMiddleware checks for a valid token in the Authorization header
func AuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		token := strings.TrimSpace(authHeader)
		if !strings.HasPrefix(token, bearerPrefix) {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token = strings.TrimPrefix(token, bearerPrefix)
		if !isValidToken(token) {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	})
}

// isValidToken validates the JWT token
func isValidToken(tokenStr string) bool {
	// Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return JWTSecret, nil
	})

	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false
	}

	// Check the claims (e.g., expiration)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				fmt.Println("Token is expired")
				return false
			}
		}

		// We can add other claim checks here (e.g., issuer, subject, roles)

		return true // Token is valid
	}

	return false // Token is invalid
}
