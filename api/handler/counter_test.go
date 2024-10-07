package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"gounter/api/handler"
	"gounter/internal/model"
	"gounter/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCounter(t *testing.T) {
	mockService := new(mocks.Service)
	h := handler.NewHandler(mockService)

	testCases := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func()
		expectedStatus int
	}{
		{
			name:        "CreateCounter Success",
			requestBody: map[string]string{"name": "testCounter"},
			mockFunc: func() {
				mockService.On("CreateCounter", mock.Anything, "testCounter").
					Return(&model.Counter{ID: uuid.New(), Name: "testCounter"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "CreateCounter Invalid JSON",
			requestBody:    "invalid",
			mockFunc:       func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var body *bytes.Buffer
			if tc.requestBody != nil {
				b, _ := json.Marshal(tc.requestBody)
				body = bytes.NewBuffer(b)
			} else {
				body = &bytes.Buffer{}
			}

			req, err := http.NewRequest("POST", "/counter", body)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			tc.mockFunc()

			h.CreateCounter(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestIncrementCounter(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(*mocks.Service)
		expectedStatus int
	}{
		{
			name:        "IncrementCounter Success",
			requestBody: map[string]uuid.UUID{"id": uuid.New()},
			mockFunc: func(mockService *mocks.Service) {
				mockService.On("IncrementCounter", mock.Anything, mock.Anything).
					Return(&model.Counter{ID: uuid.New(), Name: "testCounter"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "IncrementCounter Invalid JSON",
			requestBody:    "invalid",
			mockFunc:       func(mockService *mocks.Service) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "IncrementCounter Service Error",
			requestBody: map[string]uuid.UUID{"id": uuid.New()},
			mockFunc: func(mockService *mocks.Service) {
				mockService.On("IncrementCounter", mock.Anything, mock.Anything).
					Return(nil, errors.New("Service error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(mocks.Service)
			h := handler.NewHandler(mockService)

			var body *bytes.Buffer
			if tc.requestBody != nil {
				b, _ := json.Marshal(tc.requestBody)
				body = bytes.NewBuffer(b)
			} else {
				body = &bytes.Buffer{}
			}

			req, err := http.NewRequest("POST", "/counter/increment", body)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			tc.mockFunc(mockService)

			h.IncrementCounter(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestDeleteCounter(t *testing.T) {
	testCases := []struct {
		name           string
		urlVars        map[string]string
		mockFunc       func(*mocks.Service)
		expectedStatus int
	}{
		{
			name:    "DeleteCounter Success",
			urlVars: map[string]string{"id": gofakeit.UUID()},
			mockFunc: func(mockService *mocks.Service) {
				mockService.On("SoftDeleteCounter", mock.Anything, mock.Anything).
					Return(int64(1), nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "DeleteCounter Invalid ID",
			urlVars:        map[string]string{"id": "invalid"},
			mockFunc:       func(*mocks.Service) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "DeleteCounter Service Error",
			urlVars: map[string]string{"id": gofakeit.UUID()},
			mockFunc: func(mockService *mocks.Service) {
				mockService.On("SoftDeleteCounter", mock.Anything, mock.Anything).
					Return(int64(0), errors.New("Service error"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("DELETE", "/counter/delete?id="+tc.urlVars["id"], nil)
			assert.NoError(t, err)

			req = mux.SetURLVars(req, tc.urlVars)

			mockService := new(mocks.Service)
			tc.mockFunc(mockService)

			h := handler.NewHandler(mockService)

			rr := httptest.NewRecorder()
			h.DeleteCounter(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockService.AssertExpectations(t)
		})
	}
}
