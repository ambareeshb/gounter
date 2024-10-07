package model

import (
	"time"

	"github.com/google/uuid"
)

// Counter represents the counter model
type Counter struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	Value     int64      `db:"value"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"` // Pointer to support nullable
}
