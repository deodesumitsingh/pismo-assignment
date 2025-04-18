package v1

import (
	"errors"
	"net/http"

	"github.com/deodesumitsingh/pismo/internal/api/types/req"
	"github.com/deodesumitsingh/pismo/internal/api/types/res"
	"github.com/deodesumitsingh/pismo/internal/api/util"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/deodesumitsingh/pismo/internal/service"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	*Handler[service.TransactionService]
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		Handler: NewHandler(service),
	}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	var r req.TransactionReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, res.NewResp(nil, util.ErrorParser(err)))
		return
	}

	transaction, err := h.service.Create(r)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, repository.ErrAccountDosentExists):
			fallthrough
		case errors.Is(err, repository.ErrOperationNotSupported):
			status = http.StatusNotFound
		case errors.Is(err, service.ErrAmountAndOperationTypeMismatch):
			status = http.StatusConflict
		}

		c.JSON(status, res.NewResp(nil, err))
		return
	}

	c.JSON(http.StatusCreated, res.NewResp(transaction, err))
	return
}
