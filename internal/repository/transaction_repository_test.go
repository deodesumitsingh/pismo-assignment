package repository_test

import (
	"testing"

	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/stretchr/testify/assert"
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
		account,
		operationType,
		100.0,
	)

	assert.Equal(t, data.AccountId, account.ID)
	assert.Equal(t, data.Amount, 100.0)
	assert.Equal(t, data.ID, 1)
	assert.Equal(t, data.OperationTypeId, operationType.ID)
	assert.NoError(t, err)
}
