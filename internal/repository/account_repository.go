package repository

import (
	"errors"

	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
)

var (
	ErrAccountNumberExists = errors.New("Account with number already present")
	ErrAccountDosentExists = errors.New("Account is not present")
)

type AccountRepository interface {
	Create(string) (model.Account, error)
	AccountById(int) (model.Account, error)
}

type AccountRepositoryMem struct {
	data map[int]model.Account
}

func NewAccountRepository(c *config.AppConfig) AccountRepository {
	return &AccountRepositoryMem{
		data: make(map[int]model.Account),
	}
}

func (r *AccountRepositoryMem) Create(number string) (model.Account, error) {
	for id := range r.data {
		if number == r.data[id].Number {
			return model.Account{}, ErrAccountNumberExists
		}
	}

	account := model.Account{
		ID:     len(r.data) + 1,
		Number: number,
	}

	r.data[account.ID] = account

	return account, nil
}

func (a *AccountRepositoryMem) AccountById(id int) (model.Account, error) {
	account, ok := a.data[id]
	if !ok {
		return model.Account{}, ErrAccountDosentExists
	}

	return account, nil
}
