package service

import (
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
)

type AccountService interface {
	Create(string) (model.Account, error)
	AccountById(int) (model.Account, error)
}

type AccountServiceImpl struct {
	db repository.AccountRepository
}

func NewAccountService(repository repository.AccountRepository) AccountService {
	return &AccountServiceImpl{
		db: repository,
	}
}

func (s *AccountServiceImpl) Create(number string) (model.Account, error) {
	return s.db.Create(number)
}

func (s *AccountServiceImpl) AccountById(accountId int) (model.Account, error) {
	return s.db.AccountById(accountId)
}
