package repository

import (
	"sync"
	"time"

	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
)

type TransactionRepository interface {
	Save(account model.Account, operationType model.OperationType, amount float64) (model.Transaction, error)
}

type TransactionRepositoryMem struct {
	data map[int]model.Transaction
	mu   sync.Mutex
}

func NewTransactionRepository(c *config.AppConfig) TransactionRepository {
	return &TransactionRepositoryMem{
		data: make(map[int]model.Transaction),
	}
}

func (t *TransactionRepositoryMem) Save(account model.Account, operationType model.OperationType, amount float64) (model.Transaction, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	transaction := model.Transaction{
		ID:              len(t.data) + 1,
		AccountId:       account.ID,
		OperationTypeId: operationType.ID,
		Amount:          amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	t.data[transaction.ID] = transaction

	return transaction, nil
}
