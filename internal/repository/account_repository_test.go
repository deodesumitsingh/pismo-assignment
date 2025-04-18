package repository_test

import (
	"strconv"
	"testing"

	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountRepositoryMem_Create(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(r repository.AccountRepository)
		input         string
		expectedData  model.Account
		expectedError error
	}{
		{
			name: "failure_account_number_already_exists",
			setup: func(r repository.AccountRepository) {
				r.Create("1234")
			},
			input:         "1234",
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountNumberExists,
		},
		{
			name:  "sucess_account_created",
			setup: func(r repository.AccountRepository) {},
			input: "1234",
			expectedData: model.Account{
				ID:     1,
				Number: "1234",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accountRepositoryMem := repository.NewAccountRepository(&config.AppConfig{})
			tt.setup(accountRepositoryMem)

			data, err := accountRepositoryMem.Create(tt.input)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAccountRepositoryMem_AccountById(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(r repository.AccountRepository)
		input         string
		expectedData  model.Account
		expectedError error
	}{
		{
			name:          "failure_account_number_doesnt_exists",
			setup:         func(r repository.AccountRepository) {},
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountDosentExists,
			input:         "1",
		},
		{
			name: "sucess_gets_account_details",
			setup: func(r repository.AccountRepository) {
				r.Create("12345")
			},
			expectedData: model.Account{
				ID:     1,
				Number: "12345",
			},
			expectedError: nil,
			input:         "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accoutRepository := repository.NewAccountRepository(&config.AppConfig{})

			tt.setup(accoutRepository)

			input, err := strconv.Atoi(tt.input)
			require.NoError(t, err)

			data, err := accoutRepository.AccountById(input)

			assert.Equal(t, data, tt.expectedData)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
