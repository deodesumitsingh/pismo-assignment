package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/deodesumitsingh/pismo/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckHealtch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/healthcheck", handler.HealthCheck)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)

	router.ServeHTTP(w, r)

	assert.Equal(t, w.Code, http.StatusOK)
}
