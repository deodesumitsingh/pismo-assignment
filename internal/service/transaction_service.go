package service

import (
	"errors"

	"github.com/deodesumitsingh/pismo/internal/api/types/req"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
)

var (
	ErrAmountAndOperationTypeMismatch = errors.New("Amount and OperationType is not matching")
)

type TransactionService interface {
	Create(req.TransactionReq) (model.Transaction, error)
}

type TransactionServiceImpl struct {
	accountDb       repository.AccountRepository
	operationTypeDb repository.OperationRepository
	db              repository.TransactionRepository
}

type TransactionServiceParam struct {
	AccountRepo       repository.AccountRepository
	OperationTypeRepo repository.OperationRepository
	TransctionRepo    repository.TransactionRepository
}

func NewTransactionService(p TransactionServiceParam) TransactionService {
	return &TransactionServiceImpl{
		db:              p.TransctionRepo,
		accountDb:       p.AccountRepo,
		operationTypeDb: p.OperationTypeRepo,
	}
}

func (s *TransactionServiceImpl) Create(req req.TransactionReq) (model.Transaction, error) {
	account, err := s.accountDb.AccountById(req.AccountID)
	if err != nil {
		return model.Transaction{}, err
	}

	operationType, err := s.operationTypeDb.OperationTypeById(req.OperationTypeId)
	if err != nil {
		return model.Transaction{}, err
	}

	if !operationType.ValidOperation(req.Amount) {
		return model.Transaction{}, ErrAmountAndOperationTypeMismatch
	}

	return s.db.Save(account, operationType, req.Amount)
}
