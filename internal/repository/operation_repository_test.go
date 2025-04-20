package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOperationRepository_OperationTypeById(t *testing.T) {
	operationRepository := repository.NewOperationRepository(&config.AppConfig{})

	tests := []struct {
		name          string
		input         int
		expectedData  model.OperationType
		expectedError error
	}{
		{
			name:          "failure_operation_type_is_not_present",
			input:         10,
			expectedData:  model.OperationType{},
			expectedError: repository.ErrOperationNotSupported,
		},
		{
			name:  "success_returns_operation_type",
			input: 1,
			expectedData: model.OperationType{
				ID:          1,
				Description: "Normal Purchase",
				Mode:        model.DEBIT,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := operationRepository.OperationTypeById(tt.input)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestOperationRepositoryDb_OperationTypeById(t *testing.T) {
	query := `-- name: GetOperationById :one
SELECT id, description, mode, created_at, updated_at FROM operation_types WHERE id = $1
`

	tests := []struct {
		name          string
		input         int
		expectedData  model.OperationType
		expectedError error
		mockService   func(m sqlmock.Sqlmock)
	}{
		{
			name:  "failure_operation_type_is_not_present",
			input: 11,
			mockService: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(query).
					WillReturnError(sql.ErrNoRows)
			},
			expectedData:  model.OperationType{},
			expectedError: repository.ErrOperationNotSupported,
		},
		{
			name:  "success_returns_operation_type",
			input: 1,
			mockService: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "description", "mode", "created_at", "updated_at"}).AddRow(1, "Test", "DEBIT", time.Now(), time.Now())
				m.ExpectQuery(query).WillReturnRows(rows)
			},
			expectedData: model.OperationType{
				ID:          1,
				Description: "Test",
				Mode:        model.DEBIT,
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)

			r := repository.NewOperationRepository(&config.AppConfig{Db: db})

			tt.mockService(mock)

			data, err := r.OperationTypeById(tt.input)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
