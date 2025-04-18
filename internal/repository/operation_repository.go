package repository

import (
	"errors"

	"github.com/deodesumitsingh/pismo/config"
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
	return &OperationRepositoryMem{
		data: seedDataForInMemOperation(),
	}
}

func (o *OperationRepositoryMem) OperationTypeById(id int) (model.OperationType, error) {
	operation, ok := o.data[id]
	if !ok {
		return model.OperationType{}, ErrOperationNotSupported
	}

	return operation, nil
}
