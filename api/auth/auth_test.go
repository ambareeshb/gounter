package auth_test

import (
	"fmt"
	"gounter/api/auth"
	"gounter/util"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access granted"))
	})

	validJWT, err := util.GenerateValidJWT()
	assert.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid token",
			authHeader:     fmt.Sprintf("Bearer %s", validJWT),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing Authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid token format",
			authHeader:     "InvalidTokenFormat",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Empty Bearer token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Expired or Invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			handler := auth.AuthorizationMiddleware(testHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
