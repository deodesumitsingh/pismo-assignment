package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionRepository_Create(t *testing.T) {
	transactionRepositoryMem := repository.NewTransactionRepository(&config.AppConfig{})
	account := model.Account{
		ID:     1,
		Number: "1234",
	}
	operationType := model.OperationType{
		ID:          1,
		Description: "Testing",
		Mode:        model.DEBIT,
	}

	data, err := transactionRepositoryMem.Save(
		1,
		1,
		100.0,
	)

	assert.Equal(t, data.AccountId, account.ID)
	assert.Equal(t, data.Amount, 100.0)
	assert.Equal(t, data.ID, 1)
	assert.Equal(t, data.OperationTypeId, operationType.ID)
	assert.NoError(t, err)
}

func TestTransactionRepostiryDb_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)

	query := `-- name: CreateTransaction :one
INSERT INTO transactions(account_id, operation_type_id, amount) VALUES($1, $2, $3)
RETURNING id, amount, account_id, operation_type_id, created_at, updated_at
`
	db_time := time.Now()
	rows := sqlmock.NewRows([]string{"id", "amount", "account_id", "operation_type_id", "created_at", "updated_at"}).
		AddRow(1, 100.00, 1, 2, db_time, db_time)
	mock.ExpectQuery(query).WithArgs(1, 2, 100.00).WillReturnRows(rows)

	r := repository.NewTransactionRepository(&config.AppConfig{Db: db})

	data, err := r.Save(1, 2, 100.00)

	assert.NoError(t, err)
	assert.Equal(t, model.Transaction{
		ID:              1,
		AccountId:       1,
		OperationTypeId: 2,
		Amount:          100.00,
		CreatedAt:       db_time,
		UpdatedAt:       db_time,
	}, data)

	assert.NoError(t, mock.ExpectationsWereMet())
}
