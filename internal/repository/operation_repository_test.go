package repository_test

import (
	"testing"

	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestOperationRepository_OperationTypeById(t *testing.T) {
	operationRepository := repository.NewOperationRepository(&config.AppConfig{})

	tests := []struct {
		name          string
		input         int
		expectedData  model.OperationType
		expectedError error
	}{
		{
			name:          "failure_operation_type_is_not_present",
			input:         10,
			expectedData:  model.OperationType{},
			expectedError: repository.ErrOperationNotSupported,
		},
		{
			name:  "success_returns_operation_type",
			input: 1,
			expectedData: model.OperationType{
				ID:          1,
				Description: "Normal Purchase",
				Mode:        model.DEBIT,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := operationRepository.OperationTypeById(tt.input)

			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}
