package main

import (
	"github.com/deodesumitsingh/pismo/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/healthcheck", handler.HealthCheck)

	router.Run()
}
