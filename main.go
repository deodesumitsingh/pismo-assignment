package main

import (
	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	c := config.NewAppConfig()

	router := gin.Default()

	router.GET("/healthcheck", handler.HealthCheck)

	router.Run(c.ListnerAddr())
}
