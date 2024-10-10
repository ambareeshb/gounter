package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

const serverURL = "http://localhost:8081"
const jwtSecret = "my-secret-key" // Use the same secret key as in your app

type CounterResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestCounterIntegration(t *testing.T) {
	// Generate a valid JWT token to use for he requests
	token, err := generateValidJWT()
	assert.NoError(t, err)

	// Step 1: Create a new counter
	createCounterBody := map[string]string{"name": "testCounter"}
	createCounterResponse := sendRequest(t, "POST", "/counter/create", createCounterBody, http.StatusCreated, token)

	var createResponse CounterResponse
	err = json.Unmarshal(createCounterResponse, &createResponse)
	assert.NotEmpty(t, createResponse.ID)
	assert.NoError(t, err)

	counterID := createResponse.ID
	assert.NotEmpty(t, counterID)

	// Step 2: Increment the counter twice
	incrementCounterBody := map[string]string{"id": counterID}
	sendRequest(t, "POST", "/counter/increment", incrementCounterBody, http.StatusOK, token)
	sendRequest(t, "POST", "/counter/increment", incrementCounterBody, http.StatusOK, token)

	// Step 3: Delete the counter
	deleteCounterURL := "/counter/delete?id=" + counterID
	sendRequest(t, "DELETE", deleteCounterURL, nil, http.StatusOK, token)
}

// Helper function to send requests with Bearer token
func sendRequest(t *testing.T, method string, endpoint string, body interface{}, expectedStatus int, token string) []byte {
	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		assert.NoError(t, err)
	}

	req, err := http.NewRequest(method, serverURL+endpoint, bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("Authorization", "Bearer "+token)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, expectedStatus, resp.StatusCode)

	return responseBody
}

// Helper function to generate a valid JWT token
func generateValidJWT() (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(time.Minute * 5).Unix(), // Set expiration to 5 minutes from now
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
