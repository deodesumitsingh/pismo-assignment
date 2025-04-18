package v1_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/deodesumitsingh/pismo/internal/api/handler/v1"
	"github.com/deodesumitsingh/pismo/internal/api/types/req"
	"github.com/deodesumitsingh/pismo/internal/api/types/res"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/deodesumitsingh/pismo/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockTransactionHandler struct {
	mock.Mock
}

func (m *mockTransactionHandler) Create(r req.TransactionReq) (model.Transaction, error) {
	args := m.Called(r)
	return args.Get(0).(model.Transaction), args.Error(1)
}

func TestTransactionHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		mockService    func(m *mockTransactionHandler)
		mockVerify     func(m *mockTransactionHandler)
		input          string
		expectedBody   res.Resp
	}{
		{
			name:           "failure_invalid_account_id_type",
			expectedStatus: http.StatusBadRequest,
			mockService:    func(m *mockTransactionHandler) {},
			mockVerify: func(m *mockTransactionHandler) {
				m.AssertNotCalled(t, "Create")
			},
			input:        "{\"account_id\": \"\", \"operation_type_id\": 1, \"amount\": 100}",
			expectedBody: res.NewResp(nil, errors.New("account_id field has invalid data type")),
		},
		{
			name:           "failure_account_is_not_present",
			expectedStatus: http.StatusNotFound,
			mockService: func(m *mockTransactionHandler) {
				m.On("Create", req.TransactionReq{
					AccountID:       1,
					OperationTypeId: 1,
					Amount:          100,
				}).Return(model.Transaction{}, repository.ErrAccountDosentExists)
			},
			input:        "{\"account_id\": 1, \"operation_type_id\": 1, \"amount\": 100}",
			expectedBody: res.NewResp(nil, repository.ErrAccountDosentExists),
		},
		{
			name:           "failure_invalid_operation_type_id",
			expectedStatus: http.StatusBadRequest,
			mockVerify: func(m *mockTransactionHandler) {
				m.AssertNotCalled(t, "Create")
			},
			mockService:  func(m *mockTransactionHandler) {},
			input:        "{\"account_id\": 1, \"operation_type_id\": \"1\", \"amount\": 100}",
			expectedBody: res.NewResp(nil, errors.New("operation_type_id field has invalid data type")),
		},
		{
			name:           "failure_invalid_operation_type_is_not_present",
			expectedStatus: http.StatusNotFound,
			mockService: func(m *mockTransactionHandler) {
				m.On("Create", req.TransactionReq{
					AccountID:       1,
					OperationTypeId: 1,
					Amount:          100,
				}).Return(model.Transaction{}, repository.ErrOperationNotSupported)
			},
			input:        "{\"account_id\": 1, \"operation_type_id\": 1, \"amount\": 100}",
			expectedBody: res.NewResp(nil, repository.ErrOperationNotSupported),
		},
		{
			name:           "failure_operation_type_and_flow_of_money_is_not_supported",
			expectedStatus: http.StatusConflict,
			mockService: func(m *mockTransactionHandler) {
				m.On("Create", req.TransactionReq{
					AccountID:       1,
					OperationTypeId: 1,
					Amount:          100,
				}).Return(model.Transaction{}, service.ErrAmountAndOperationTypeMismatch)
			},
			input:        "{\"account_id\": 1, \"operation_type_id\": 1, \"amount\": 100}",
			expectedBody: res.NewResp(nil, service.ErrAmountAndOperationTypeMismatch),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.Default()

			mock := new(mockTransactionHandler)
			handler := v1.NewTransactionHandler(mock)
			tt.mockService(mock)

			router.POST("/transactions", handler.Create)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte(tt.input)))
			require.NoError(t, err)

			router.ServeHTTP(w, req)

			data, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, string(data), w.Body.String())

			if tt.mockVerify != nil {
				tt.mockVerify(mock)
			}
		})
	}
}
