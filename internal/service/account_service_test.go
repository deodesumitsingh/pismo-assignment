package service_test

import (
	"strconv"
	"testing"

	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/deodesumitsingh/pismo/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAccountRepository struct {
	mock.Mock
}

func (m *mockAccountRepository) Create(number string) (model.Account, error) {
	args := m.Called(number)
	return args.Get(0).(model.Account), args.Error(1)
}

func (m *mockAccountRepository) AccountById(accountId int) (model.Account, error) {
	args := m.Called(accountId)
	return args.Get(0).(model.Account), args.Error(1)
}

type test struct {
	name          string
	mockService   func(m *mockAccountRepository)
	input         string
	expectedData  model.Account
	expectedError error
}

func TestAccountService_Create(t *testing.T) {
	account := model.Account{
		ID:     1,
		Number: "12345",
	}

	tests := []test{
		{
			name: "failure_accout_number_already_exits",
			mockService: func(m *mockAccountRepository) {
				m.On("Create", "12345").Return(model.Account{}, repository.ErrAccountNumberExists)
			},
			input:         "12345",
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountNumberExists,
		},
		{
			name: "success_account_created",
			mockService: func(m *mockAccountRepository) {
				m.On("Create", "12345").Return(account, nil)
			},
			input:         "12345",
			expectedData:  account,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAccountRepository{}

			accountService := service.NewAccountService(mock)
			tt.mockService(mock)

			data, err := accountService.Create(tt.input)

			assert.Equal(t, data, tt.expectedData)
			assert.Equal(t, err, tt.expectedError)
			mock.AssertExpectations(t)
		})
	}
}

func TestAccountService_AccountById(t *testing.T) {
	account := model.Account{
		ID:     1,
		Number: "12345",
	}

	tests := []test{
		{
			name: "failure_account_number_dosent_exits",
			mockService: func(m *mockAccountRepository) {
				m.On("AccountById", 10).Return(model.Account{}, repository.ErrAccountDosentExists)
			},
			input:         "10",
			expectedData:  model.Account{},
			expectedError: repository.ErrAccountDosentExists,
		},
		{
			name: "success_gets_the_account",
			mockService: func(m *mockAccountRepository) {
				m.On("AccountById", 1).Return(account, nil)
			},
			input:         "1",
			expectedData:  account,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAccountRepository{}

			accountService := service.NewAccountService(mock)
			tt.mockService(mock)

			i, err := strconv.Atoi(tt.input)
			require.NoError(t, err)
			data, err := accountService.AccountById(i)

			assert.Equal(t, data, tt.expectedData)
			assert.Equal(t, err, tt.expectedError)
			mock.AssertExpectations(t)
		})
	}
}
