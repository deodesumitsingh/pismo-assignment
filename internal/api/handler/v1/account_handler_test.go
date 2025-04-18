package v1_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/deodesumitsingh/pismo/internal/api/handler/v1"
	"github.com/deodesumitsingh/pismo/internal/api/types/res"
	"github.com/deodesumitsingh/pismo/internal/model"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockAccountService struct {
	mock.Mock
}

func (m *mockAccountService) Create(number string) (model.Account, error) {
	args := m.Called(number)
	return args.Get(0).(model.Account), args.Error(1)
}

func (m *mockAccountService) AccountById(accountId int) (model.Account, error) {
	args := m.Called(accountId)
	return args.Get(0).(model.Account), args.Error(1)
}

func TestAccountHandler_Create(t *testing.T) {
	account := model.Account{ID: 1, Number: "12345"}

	tests := []struct {
		name           string
		input          string
		mockService    func(mock *mockAccountService)
		expectedStatus int
		expectedBody   res.Resp
		verifyMock     func(mock *mockAccountService)
	}{
		{
			name:  "failure",
			input: "{\"document_number\": \"\"}",
			mockService: func(mock *mockAccountService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   res.NewResp(nil, errors.New("Number field is required")),
			verifyMock: func(mock *mockAccountService) {
				mock.AssertNotCalled(t, "Create")
			},
		},
		{
			name:  "success",
			input: "{\"document_number\": \"12345\"}",
			mockService: func(mock *mockAccountService) {
				mock.On("Create", "12345").Return(account, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   res.NewResp(account, nil),
		},
		{
			name:  "document number already exists",
			input: "{\"document_number\": \"12345\"}",
			mockService: func(mock *mockAccountService) {
				mock.On("Create", "12345").Return(model.Account{}, repository.ErrAccountNumberExists)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   res.NewResp(nil, repository.ErrAccountNumberExists),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.Default()

			mockService := &mockAccountService{}
			tt.mockService(mockService)

			handler := v1.NewAccountHandler(mockService)
			router.POST("/accounts", handler.Create)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(
				http.MethodPost,
				"/accounts",
				bytes.NewReader([]byte(tt.input)),
			)
			router.ServeHTTP(w, req)

			data, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, string(data), w.Body.String())
			mockService.AssertExpectations(t)

			if tt.verifyMock != nil {
				tt.verifyMock(mockService)
			}
		})
	}
}

func TestAccountHandler_GetAccount(t *testing.T) {
	tests := []struct {
		name           string
		accountId      string
		expectedStatus int
		expectedBody   res.Resp
		mockService    func(m *mockAccountService)
		mockVerify     func(m *mockAccountService)
	}{
		{
			name:           "failure_account_not_found",
			accountId:      "12345",
			expectedStatus: http.StatusNotFound,
			expectedBody:   res.NewResp(nil, repository.ErrAccountDosentExists),
			mockService: func(m *mockAccountService) {
				m.On("AccountById", 12345).Return(model.Account{}, repository.ErrAccountDosentExists)
			},
		},
		{
			name:           "success",
			accountId:      "12345",
			expectedStatus: http.StatusOK,
			expectedBody:   res.NewResp(model.Account{ID: 1, Number: "12345"}, nil),
			mockService: func(m *mockAccountService) {
				m.On("AccountById", 12345).Return(model.Account{
					ID:     1,
					Number: "12345",
				}, nil)
			},
		},
		{
			name:           "failure_account_id_invalid",
			accountId:      "invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   res.NewResp(nil, v1.ErrInvalidAccountId),
			mockService:    func(m *mockAccountService) {},
			mockVerify: func(m *mockAccountService) {
				m.AssertNotCalled(t, "AccountById")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			router := gin.Default()

			mock := new(mockAccountService)
			tt.mockService(mock)

			handler := v1.NewAccountHandler(mock)

			router.GET("/accounts/:accountId", handler.GetAccount)

			w := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, "/accounts/"+tt.accountId, nil)
			require.NoError(t, err)
			router.ServeHTTP(w, req)

			mock.AssertExpectations(t)
			assert.Equal(t, tt.expectedStatus, w.Code)

			data, err := json.Marshal(tt.expectedBody)
			require.NoError(t, err)

			assert.JSONEq(t, string(data), w.Body.String())

			if tt.mockVerify != nil {
				tt.mockVerify(mock)
			}
		})
	}
}
