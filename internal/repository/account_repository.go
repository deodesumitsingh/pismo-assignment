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

type AccountRepositoryDb struct {
	db *database.Queries
}

func NewAccountRepository(c *config.AppConfig) AccountRepository {
	if c.Db != nil {
		return &AccountRepositoryDb{
			db: database.New(c.Db),
		}
	}

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

func (r *AccountRepositoryMem) AccountById(id int) (model.Account, error) {
	account, ok := r.data[id]
	if !ok {
		return model.Account{}, ErrAccountDosentExists
	}

	return account, nil
}

func (r *AccountRepositoryDb) Create(number string) (model.Account, error) {
	q, err := r.db.GetAccountByNumber(context.Background(), number)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return model.Account{}, err
	}
	if q.Number == number {
		return model.Account{}, ErrAccountNumberExists
	}

	data, err := r.db.CreateAccount(context.Background(), number)
	if err != nil {
		return model.Account{}, err
	}

	return model.Account{
		ID:     int(data.ID),
		Number: data.Number,
	}, nil
}

func (r *AccountRepositoryDb) AccountById(id int) (model.Account, error) {
	data, err := r.db.GetAccountById(context.Background(), int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrAccountDosentExists
		}
		return model.Account{}, err
	}

	return model.Account{
		ID:     int(data.ID),
		Number: data.Number,
	}, err
}
