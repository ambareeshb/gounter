package repository

import (
	"context"
	"gounter/internal/model"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	CreateCounterSQL = `
		INSERT INTO counter (id, name, value, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, value;`

	IncrementCounterSQL = `
		UPDATE counter
		SET value = value + 1
		WHERE id = $1
		RETURNING id, name, value;`

	SoftDeleteCounter = `
		UPDATE counter 
		SET deleted_at = $1 
		WHERE id = $2 AND deleted_at IS NULL;`
)

// CounterRepository defines the interface for counter operations
type CounterRepository interface {
	CreateCounter(ctx context.Context, name string) (*model.Counter, error)
}

// Counter struct do operation counter table in db.
type Counter struct {
	db *sqlx.DB
}

// NewCounterRepository creates a new instance of CounterRepository
func New(db *sqlx.DB) *Counter {
	return &Counter{db: db}
}

// CreateCounter inserts a new counter into the database and returns the created counter
func (r *Counter) CreateCounter(ctx context.Context, name string) (*model.Counter, error) {
	now := time.Now().UTC()
	id := uuid.New()

	var counter model.Counter

	// Perform the insert and return the created counter
	err := r.db.QueryRowContext(ctx, CreateCounterSQL, id, name, 0, now, now).
		Scan(&counter.ID, &counter.Name, &counter.Value)
	if err != nil {
		return nil, err
	}

	return &counter, nil
}

// IncrementCounter increments the counter by 1 and returns the new value
func (r *Counter) IncrementCounter(ctx context.Context, id uuid.UUID) (*model.Counter, error) {
	var counter model.Counter

	err := r.db.QueryRowContext(ctx, IncrementCounterSQL, id).
		Scan(&counter.ID, &counter.Name, &counter.Value)
	if err != nil {
		return nil, err
	}

	return &counter, nil
}

// SoftDeleteCounter marks the counter as deleted by setting the deleted_at timestamp.
// It returns the number of rows affected.
func (r *Counter) SoftDeleteCounter(ctx context.Context, id uuid.UUID) (int64, error) {
	now := time.Now().UTC()

	result, err := r.db.ExecContext(ctx, SoftDeleteCounter, now, id)
	if err != nil {
		return 0, err
	}

	// Check how many rows were affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	// Return the number of affected rows
	return rowsAffected, nil
}
