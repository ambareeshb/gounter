package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const serverURL = "http://localhost:8081"

type CounterResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestCounterIntegration(t *testing.T) {
	// Step 1: Create a new counter
	createCounterBody := map[string]string{"name": "testCounter"}
	createCounterResponse := sendRequest(t, "POST", "/counter/create", createCounterBody, http.StatusCreated)

	var createResponse CounterResponse
	err := json.Unmarshal(createCounterResponse, &createResponse)
	assert.NoError(t, err)

	counterID := createResponse.ID
	assert.NotEmpty(t, counterID)

	// Step 2: Increment the counter twice
	incrementCounterBody := map[string]string{"id": counterID}
	sendRequest(t, "POST", "/counter/increment", incrementCounterBody, http.StatusOK)
	sendRequest(t, "POST", "/counter/increment", incrementCounterBody, http.StatusOK)

	// Step 3: Delete the counter
	deleteCounterURL := "/counter/delete?id=" + counterID
	sendRequest(t, "DELETE", deleteCounterURL, nil, http.StatusOK)
}

func sendRequest(t *testing.T, method string, endpoint string, body interface{}, expectedStatus int) []byte {
	var requestBody []byte
	var err error

	if body != nil {
		requestBody, err = json.Marshal(body)
		assert.NoError(t, err)
	}

	// Send the request to the server
	req, err := http.NewRequest(method, serverURL+endpoint, bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Assert the status code
	assert.Equal(t, expectedStatus, resp.StatusCode)

	return responseBody
}
