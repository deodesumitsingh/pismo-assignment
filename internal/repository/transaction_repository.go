package repository

import (
	"context"
	"sync"
	"time"

	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/database"
	"github.com/deodesumitsingh/pismo/internal/model"
)

type TransactionRepository interface {
	Save(account_id, operation_type_id int, amount float64) (model.Transaction, error)
}

type TransactionRepositoryMem struct {
	data map[int]model.Transaction
	mu   sync.Mutex
}

type TransactionRepositoryDb struct {
	db *database.Queries
}

func NewTransactionRepository(c *config.AppConfig) TransactionRepository {
	if c.Db != nil {
		return &TransactionRepositoryDb{
			db: database.New(c.Db),
		}
	}

	return &TransactionRepositoryMem{
		data: make(map[int]model.Transaction),
	}
}

func (r *TransactionRepositoryMem) Save(account_id, operation_type_id int, amount float64) (model.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	transaction := model.Transaction{
		ID:              len(r.data) + 1,
		AccountId:       account_id,
		OperationTypeId: operation_type_id,
		Amount:          amount,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	r.data[transaction.ID] = transaction

	return transaction, nil
}

func (r *TransactionRepositoryDb) Save(account_id, operation_type_id int, amount float64) (model.Transaction, error) {
	transaction, err := r.db.CreateTransaction(context.Background(), database.CreateTransactionParams{
		AccountID:       int32(account_id),
		OperationTypeID: int32(operation_type_id),
		Amount:          amount,
	})
	if err != nil {
		return model.Transaction{}, err
	}

	return model.Transaction{
		ID:              int(transaction.ID),
		AccountId:       int(transaction.AccountID),
		OperationTypeId: int(transaction.OperationTypeID),
		Amount:          transaction.Amount,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}, nil
}
