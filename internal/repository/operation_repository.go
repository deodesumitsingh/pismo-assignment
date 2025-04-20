package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/database"
	"github.com/deodesumitsingh/pismo/internal/model"
)

var (
	ErrOperationNotSupported = errors.New("Operation is not supported")
)

type OperationRepository interface {
	OperationTypeById(id int) (model.OperationType, error)
}

type OperationRepositoryMem struct {
	data map[int]model.OperationType
}

type OperationRepositoryDb struct {
	db *database.Queries
}

func seedDataForInMemOperation() map[int]model.OperationType {
	operations := []model.OperationType{
		{ID: 1, Description: "Normal Purchase", Mode: model.DEBIT},
		{ID: 2, Description: "Purchase With Installments", Mode: model.DEBIT},
		{ID: 3, Description: "Withdrawal", Mode: model.DEBIT},
		{ID: 4, Description: "Credit Voucher", Mode: model.CREDIT},
	}

	d := make(map[int]model.OperationType)

	for _, operation := range operations {
		d[operation.ID] = operation
	}

	return d
}

func NewOperationRepository(c *config.AppConfig) OperationRepository {
	if c.Db != nil {
		return &OperationRepositoryDb{
			db: database.New(c.Db),
		}
	}
	return &OperationRepositoryMem{
		data: seedDataForInMemOperation(),
	}
}

func (r *OperationRepositoryMem) OperationTypeById(id int) (model.OperationType, error) {
	operation, ok := r.data[id]
	if !ok {
		return model.OperationType{}, ErrOperationNotSupported
	}

	return operation, nil
}

func (r *OperationRepositoryDb) OperationTypeById(id int) (model.OperationType, error) {
	operation, err := r.db.GetOperationById(context.Background(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrOperationNotSupported
		}
		return model.OperationType{}, err
	}

	return model.OperationType{
		ID:          int(operation.ID),
		Description: operation.Description,
		Mode:        model.Mode(operation.Mode),
	}, nil
}
