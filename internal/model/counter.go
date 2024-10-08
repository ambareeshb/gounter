package model

import (
	"time"

	"github.com/google/uuid"
)

// Counter represents the counter model
type Counter struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Value     int64      `db:"value" json:"value"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"` // Pointer to support nullable
}
