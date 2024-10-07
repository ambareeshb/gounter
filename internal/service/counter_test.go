package service_test

import (
	"context"
	"database/sql"
	"errors"
	"gounter/internal/model"
	"gounter/internal/service"
	"gounter/test/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCounterServiceCreateCounter(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(repo *mocks.Repository)
		inputName     string
		expectedError error
	}{
		{
			name: "successfully creates a counter",
			setupMock: func(repo *mocks.Repository) {
				repo.On("CreateCounter", mock.Anything, "test_counter").Return(&model.Counter{}, nil)
			},
			inputName:     "test_counter",
			expectedError: nil,
		},
		{
			name: "failed to create counter due to database error",
			setupMock: func(repo *mocks.Repository) {
				repo.On("CreateCounter", mock.Anything, "test_counter").Return(nil, errors.New("db error"))
			},
			inputName:     "test_counter",
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.Repository)
			tt.setupMock(repo)

			svc := service.NewCounterService(repo)

			ctx := context.TODO()
			_, err := svc.CreateCounter(ctx, tt.inputName)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestCounterServiceIncrementCounter(t *testing.T) {
	tests := []struct {
		name          string
		expectedValue *model.Counter
		setupMock     func(repo *mocks.Repository, id uuid.UUID)
		inputID       uuid.UUID
		expectedError error
	}{
		{
			name: "successfully increments counter",
			setupMock: func(repo *mocks.Repository, id uuid.UUID) {
				repo.On("IncrementCounter", mock.Anything, id).Return(&model.Counter{Name: "Test Counter", Value: 1}, nil)
			},
			inputID:       uuid.New(),
			expectedError: nil,
			expectedValue: &model.Counter{Name: "Test Counter", Value: 1},
		},
		{
			name: "counter not found",
			setupMock: func(repo *mocks.Repository, id uuid.UUID) {
				repo.On("IncrementCounter", mock.Anything, id).Return(nil, sql.ErrNoRows)
			},
			inputID:       uuid.New(),
			expectedError: service.ErrCounterNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.Repository)
			tt.setupMock(repo, tt.inputID)

			svc := service.NewCounterService(repo)

			ctx := context.TODO()
			newValue, err := svc.IncrementCounter(ctx, tt.inputID)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedValue, newValue)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestCounterServiceSoftDeleteCounter(t *testing.T) {
	tests := []struct {
		name             string
		setupMock        func(repo *mocks.Repository, id uuid.UUID)
		inputID          uuid.UUID
		expectedAffected int64
		expectedError    error
	}{
		{
			name: "successfully soft deletes a counter",
			setupMock: func(repo *mocks.Repository, id uuid.UUID) {
				repo.On("SoftDeleteCounter", mock.Anything, id).Return(int64(1), nil)
			},
			inputID:          uuid.New(),
			expectedAffected: 1,
			expectedError:    nil,
		},
		{
			name: "not found",
			setupMock: func(repo *mocks.Repository, id uuid.UUID) {
				repo.On("SoftDeleteCounter", mock.Anything, id).Return(int64(1), sql.ErrNoRows)
			},
			inputID:          uuid.New(),
			expectedAffected: 0,
			expectedError:    service.ErrCounterNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.Repository)
			tt.setupMock(repo, tt.inputID)

			svc := service.NewCounterService(repo)

			ctx := context.TODO()
			actualAffectedRows, err := svc.SoftDeleteCounter(ctx, tt.inputID)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedAffected, actualAffectedRows)
			}

			repo.AssertExpectations(t)
		})
	}
}
