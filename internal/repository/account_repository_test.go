package repository_test

import (
	"database/sql"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepositoryMem_Create(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(r repository.AccountRepository)
		input         string
		expectedData  model.Account
		expectedError error
	}{
		{
			name: "failure_account_number_already_exists",
			setup: func(r repository.AccountRepository) {
				r.Create("1234")
			},
			input:         "1234",
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountNumberExists,
		},
		{
			name:  "sucess_account_created",
			setup: func(r repository.AccountRepository) {},
			input: "1234",
			expectedData: model.Account{
				ID:     1,
				Number: "1234",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepositoryMem := repository.NewAccountRepository(&config.AppConfig{})
			tt.setup(accountRepositoryMem)

			data, err := accountRepositoryMem.Create(tt.input)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAccountRepositoryMem_AccountById(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(r repository.AccountRepository)
		input         string
		expectedData  model.Account
		expectedError error
	}{
		{
			name:          "failure_account_number_doesnt_exists",
			setup:         func(r repository.AccountRepository) {},
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountDosentExists,
			input:         "1",
		},
		{
			name: "sucess_gets_account_details",
			setup: func(r repository.AccountRepository) {
				r.Create("12345")
			},
			expectedData: model.Account{
				ID:     1,
				Number: "12345",
			},
			expectedError: nil,
			input:         "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accoutRepository := repository.NewAccountRepository(&config.AppConfig{})

			tt.setup(accoutRepository)

			input, err := strconv.Atoi(tt.input)
			require.NoError(t, err)

			data, err := accoutRepository.AccountById(input)

			assert.Equal(t, data, tt.expectedData)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAccountRepositoryDb_AccountById(t *testing.T) {
	query := `-- name: GetAccountById :one
SELECT id, number, created_at, updated_at FROM accounts WHERE id = ?`

	tests := []struct {
		name          string
		mockService   func(m sqlmock.Sqlmock)
		input         int
		expectedData  model.Account
		expectedError error
	}{
		{
			name: "failure_account_number_doesnt_exists",
			mockService: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(query).
					WithArgs(1).
					WillReturnError(sql.ErrNoRows)
			},
			input:         1,
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountDosentExists,
		},
		{
			name: "sucess_gets_account_details",
			mockService: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "number", "created_at", "updated_at"}).
					AddRow(1, "12345", time.Now(), time.Now())

				m.ExpectQuery(query).
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: 1,
			expectedData: model.Account{
				ID:     1,
				Number: "12345",
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			defer db.Close()

			require.NoError(t, err)

			r := repository.NewAccountRepository(&config.AppConfig{
				Db: db,
			})

			tt.mockService(mock)
			data, err := r.AccountById(tt.input)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAccountRepositoryDb_Create(t *testing.T) {
	queryGetAccountByNumber := `-- name: GetAccountByNumber :one
SELECT id, number, created_at, updated_at FROM accounts WHERE number = $1`
	queryCreateAccount := `-- name: CreateAccount :one
INSERT INTO accounts(number) values ($1)
RETURNING id, number, created_at, updated_at`

	tests := []struct {
		name          string
		mockService   func(m sqlmock.Sqlmock)
		input         string
		expectedData  model.Account
		expectedError error
	}{
		{
			name: "failure_account_number_already_exists",
			mockService: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "number", "created_at", "updated_at"}).
					AddRow(1, "12345", time.Now(), time.Now())
				m.ExpectQuery(queryGetAccountByNumber).
					WithArgs("12345").
					WillReturnRows(rows)
			},
			input:         "12345",
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountNumberExists,
		},
		{
			name: "success_account_created",
			mockService: func(m sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "number", "created_at", "updated_at"}).
					AddRow(1, "12345", time.Now(), time.Now())
				m.ExpectQuery(queryGetAccountByNumber).
					WithArgs("12345").
					WillReturnError(sql.ErrNoRows)
				m.ExpectQuery(queryCreateAccount).
					WithArgs("12345").
					WillReturnRows(rows)
			},
			input: "12345",
			expectedData: model.Account{
				ID:     1,
				Number: "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			require.NoError(t, err)
			defer db.Close()

			r := repository.NewAccountRepository(&config.AppConfig{Db: db})
			tt.mockService(mock)

			data, err := r.Create(tt.input)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedError, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
