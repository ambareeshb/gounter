package service

import (
	"context"
	"database/sql"
	"errors"
	"gounter/internal/model"

	"github.com/google/uuid"
)

// Repository defines the interface for the counter repository
type Repository interface {
	SoftDeleteCounter(ctx context.Context, id uuid.UUID) (int64, error)
	IncrementCounter(ctx context.Context, id uuid.UUID) (*model.Counter, error)
	CreateCounter(ctx context.Context, name string) (*model.Counter, error)
}

// ErrCounterNotFound is returned when a counter is not found
var ErrCounterNotFound = errors.New("counter not found")

// CounterService is an implementation of the Service interface
type CounterService struct {
	repo Repository
}

// NewCounterService creates a new instance of the counter service
func NewCounterService(repo Repository) *CounterService {
	return &CounterService{repo: repo}
}

// CreateCounter calls the repository to create a counter and returns the created counter
func (s *CounterService) CreateCounter(ctx context.Context, name string) (*model.Counter, error) {
	counter, err := s.repo.CreateCounter(ctx, name)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

// IncrementCounter increments the counter value and returns the updated counter
func (s *CounterService) IncrementCounter(ctx context.Context, id uuid.UUID) (*model.Counter, error) {
	newCounterValue, err := s.repo.IncrementCounter(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCounterNotFound
		}

		return nil, err
	}

	return newCounterValue, nil
}

// SoftDeleteCounter soft deletes a counter and returns meaningful error if the counter is already deleted or not found
func (s *CounterService) SoftDeleteCounter(ctx context.Context, id uuid.UUID) (int64, error) {
	rowsAffected, err := s.repo.SoftDeleteCounter(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrCounterNotFound
		}
		return 0, err
	}

	return rowsAffected, nil
}
