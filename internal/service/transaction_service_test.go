package service_test

import (
	"testing"
	"time"

	"github.com/deodesumitsingh/pismo/internal/api/types/req"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/deodesumitsingh/pismo/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(s string) (model.Account, error) {
	return model.Account{}, nil
}
func (m *mockRepository) AccountById(id int) (model.Account, error) {
	args := m.Called(id)
	return args.Get(0).(model.Account), args.Error(1)
}

func (m *mockRepository) OperationTypeById(id int) (model.OperationType, error) {
	args := m.Called(id)
	return args.Get(0).(model.OperationType), args.Error(1)
}

func (m *mockRepository) Save(a, o int, amount float64) (model.Transaction, error) {
	args := m.Called(a, o, amount)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func TestTransactionService_Save(t *testing.T) {
	inputCredit := req.TransactionReq{
		AccountID:       1,
		OperationTypeId: 1,
		Amount:          100.0,
	}
	inputDebit := req.TransactionReq{
		AccountID:       1,
		OperationTypeId: 1,
		Amount:          -100.0,
	}

	account := model.Account{
		ID:     1,
		Number: "12345",
	}
	operationType := model.OperationType{
		ID:          1,
		Description: "Purchase",
		Mode:        model.DEBIT,
	}
	transaction := model.Transaction{
		ID:              1,
		AccountId:       1,
		OperationTypeId: 1,
		Amount:          -100.0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	tests := []struct {
		name          string
		mockService   func(m *mockRepository)
		mockVerify    func(m *mockRepository)
		expectedData  model.Transaction
		expectedError error
		input         req.TransactionReq
	}{
		{
			name: "failure_account_id_dosen't_exists",
			mockService: func(m *mockRepository) {
				m.On("AccountById", 1).Return(model.Account{}, repository.ErrAccountDosentExists)
			},
			mockVerify: func(m *mockRepository) {
				m.AssertNotCalled(t, "OperationTypeById")
				m.AssertNotCalled(t, "Save")
			},
			expectedData:  model.Transaction{},
			expectedError: repository.ErrAccountDosentExists,
			input:         inputCredit,
		},
		{
			name: "failure_operation_type_dosen't_exists",
			mockService: func(m *mockRepository) {
				m.On("AccountById", 1).Return(account, nil)
				m.On("OperationTypeById", 1).Return(model.OperationType{}, repository.ErrOperationNotSupported)
			},
			mockVerify: func(m *mockRepository) {
				m.AssertNotCalled(t, "Save")
			},
			expectedData:  model.Transaction{},
			expectedError: repository.ErrOperationNotSupported,
			input:         inputCredit,
		},
		{
			name: "failure_operation_type_amount_mismtach",
			mockService: func(m *mockRepository) {
				m.On("AccountById", 1).Return(account, nil)
				m.On("OperationTypeById", 1).Return(operationType, nil)
			},
			mockVerify: func(m *mockRepository) {
				m.AssertNotCalled(t, "Save")
			},
			expectedData:  model.Transaction{},
			input:         inputCredit,
			expectedError: service.ErrAmountAndOperationTypeMismatch,
		},
		{
			name: "sucess_transaction_is_created",
			mockService: func(m *mockRepository) {
				m.On("AccountById", 1).Return(account, nil)
				m.On("OperationTypeById", 1).Return(operationType, nil)
				m.On("Save", account.ID, operationType.ID, inputDebit.Amount).Return(transaction, nil)
			},
			expectedData:  transaction,
			input:         inputDebit,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockRepository{}

			service := service.NewTransactionService(service.TransactionServiceParam{
				AccountRepo:       mock,
				OperationTypeRepo: mock,
				TransctionRepo:    mock,
			})

			tt.mockService(mock)

			data, err := service.Create(tt.input)

			assert.Equal(t, data, tt.expectedData)
			assert.Equal(t, err, tt.expectedError)
			mock.AssertExpectations(t)

			if tt.mockVerify != nil {
				tt.mockVerify(mock)
			}
		})
	}
}
