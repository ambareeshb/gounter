package repository_test

import (
	"context"
	"database/sql"
	"gounter/internal/model"
	counterRepository "gounter/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func TestRepositoryCreateCounter(t *testing.T) {
	tests := []struct {
		name            string
		setupMock       func(mock sqlmock.Sqlmock)
		inputName       string
		expectedCounter *model.Counter
		expectedError   error
	}{
		{
			name: "successfully creates counter",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock the insertion and return the created counter
				mock.ExpectQuery(`INSERT INTO counter \(id, name, value, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id, name, value;`).
					WithArgs(sqlmock.AnyArg(), "Test Counter", 0, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value"}).
						AddRow(gofakeit.UUID(), "Test Counter", 0))
			},
			inputName: "Test Counter",
			expectedCounter: &model.Counter{
				Name:  "Test Counter",
				Value: 0,
			},
			expectedError: nil,
		},
		{
			name: "database error during insertion",
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock a database error
				mock.ExpectQuery(`INSERT INTO counter \(id, name, value, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id, name, value;`).
					WillReturnError(sql.ErrConnDone)
			},
			inputName:       "Test Counter",
			expectedCounter: nil,
			expectedError:   sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "postgres")
			repo := counterRepository.New(sqlxDB)

			tt.setupMock(mock)

			ctx := context.TODO()
			counter, err := repo.CreateCounter(ctx, tt.inputName)

			// Validate the results
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
				require.Nil(t, counter)
			} else {
				require.NoError(t, err)
				require.NotNil(t, counter)
				require.Equal(t, tt.expectedCounter.Name, counter.Name)
				require.Equal(t, tt.expectedCounter.Value, counter.Value)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestRepositoryIncrementCounter(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, id uuid.UUID)
		id            uuid.UUID
		expectedValue *model.Counter
		expectedError error
	}{
		{
			name: "successfully increments counter",
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`UPDATE counter SET value = value \+ 1 WHERE id = \$1 RETURNING id, name, value;`).
					WithArgs(id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value"}).
						AddRow(gofakeit.UUID(), "Test Counter", 11))
			},
			id:            uuid.New(),
			expectedValue: &model.Counter{Name: "Test Counter", Value: 11},
			expectedError: nil,
		},
		{
			name: "counter not found",
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`UPDATE counter SET value = value \+ 1 WHERE id = \$1 RETURNING id, name, value;`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			id:            uuid.New(),
			expectedError: sql.ErrNoRows,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`UPDATE counter SET value = value \+ 1 WHERE id = \$1 RETURNING id, name, value;`).
					WithArgs(id).
					WillReturnError(sql.ErrConnDone)
			},
			id:            uuid.New(),
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "postgres")
			repo := counterRepository.New(sqlxDB)

			tt.setupMock(mock, tt.id)

			// Call the IncrementCounter function
			ctx := context.TODO()
			updatedCounter, err := repo.IncrementCounter(ctx, tt.id)

			// Validate the results
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedValue.Name, updatedCounter.Name)
				require.Equal(t, tt.expectedValue.Value, updatedCounter.Value)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}

func TestRepositorySoftDeleteCounter(t *testing.T) {
	tests := []struct {
		name             string
		setupMock        func(mock sqlmock.Sqlmock, id uuid.UUID)
		id               uuid.UUID
		expectedAffected int64
		expectedError    error
	}{
		{
			name: "successfully soft deletes counter",
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				// Mock the update and return 1 row affected
				mock.ExpectExec(`UPDATE counter SET deleted_at = \$1 WHERE id = \$2 AND deleted_at IS NULL;`).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(1, 1)) // 1 row affected
			},
			id:               uuid.New(),
			expectedAffected: 1,
			expectedError:    nil,
		},
		{
			name: "counter not found",
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				// Mock the update but return 0 rows affected (no such counter)
				mock.ExpectExec(`UPDATE counter SET deleted_at = \$1 WHERE id = \$2 AND deleted_at IS NULL;`).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected
			},
			id:               uuid.New(),
			expectedAffected: 0,
			expectedError:    nil,
		},
		{
			name: "counter already deleted",
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				// Mock the update but return 0 rows affected (already deleted)
				mock.ExpectExec(`UPDATE counter SET deleted_at = \$1 WHERE id = \$2 AND deleted_at IS NULL;`).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(1, 0)) // 0 rows affected
			},
			id:               uuid.New(),
			expectedAffected: 0,
			expectedError:    nil,
		},
		{
			name: "database error during update",
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				// Mock a database error
				mock.ExpectExec(`UPDATE counter SET deleted_at = \$1 WHERE id = \$2 AND deleted_at IS NULL;`).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnError(sql.ErrConnDone)
			},
			id:               uuid.New(),
			expectedAffected: 0,
			expectedError:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			sqlxDB := sqlx.NewDb(db, "postgres")
			repo := counterRepository.New(sqlxDB)

			tt.setupMock(mock, tt.id)

			ctx := context.TODO()
			rowsAffected, err := repo.SoftDeleteCounter(ctx, tt.id)

			// Validate the results
			if tt.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tt.expectedError, err)
				require.Equal(t, int64(0), rowsAffected)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedAffected, rowsAffected)
			}

			err = mock.ExpectationsWereMet()
			require.NoError(t, err)
		})
	}
}
