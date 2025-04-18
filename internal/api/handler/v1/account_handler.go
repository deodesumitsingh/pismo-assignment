package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/deodesumitsingh/pismo/internal/api/types/req"
	"github.com/deodesumitsingh/pismo/internal/api/types/res"
	"github.com/deodesumitsingh/pismo/internal/api/util"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/deodesumitsingh/pismo/internal/service"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidAccountId = errors.New("Invalid accountId format")
)

type AccountHandler struct {
	*Handler[service.AccountService]
}

func NewAccountHandler(service service.AccountService) *AccountHandler {
	return &AccountHandler{
		Handler: NewHandler(service),
	}
}

func (h *AccountHandler) Create(c *gin.Context) {
	var req req.AccountReq
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, res.NewResp(nil, util.ErrorParser(err)))
		return
	}

	account, err := h.service.Create(req.Number)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, repository.ErrAccountNumberExists) {
			status = http.StatusBadRequest
		}
		c.JSON(status, res.NewResp(nil, err))
		return
	}

	c.JSON(http.StatusCreated, res.NewResp(account, nil))
	return
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id := c.Params.ByName("accountId")
	accountId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.NewResp(nil, ErrInvalidAccountId))
		return
	}

	account, err := h.service.AccountById(int(accountId))
	if err != nil {
		c.JSON(http.StatusNotFound, res.NewResp(nil, err))
		return

	}

	c.JSON(http.StatusOK, res.NewResp(account, nil))
	return
}
