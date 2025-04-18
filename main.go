package main

import (
	"github.com/deodesumitsingh/pismo/config"
	"github.com/deodesumitsingh/pismo/internal/api/handler"
	v1 "github.com/deodesumitsingh/pismo/internal/api/handler/v1"
	"github.com/deodesumitsingh/pismo/internal/repository"
	"github.com/deodesumitsingh/pismo/internal/service"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func applicationRepositories(c *config.AppConfig) repository.ApplicationRepositories {
	return repository.ApplicationRepositories{
		AccountRepository:     repository.NewAccountRepository(c),
		OperationRepository:   repository.NewOperationRepository(c),
		TransactionRepository: repository.NewTransactionRepository(c),
	}
}

func applicationServices(repositories repository.ApplicationRepositories) service.ApplicationServices {
	return service.ApplicationServices{
		AccountService: service.NewAccountService(repositories.AccountRepository),
		TransactionService: service.NewTransactionService(service.TransactionServiceParam{
			AccountRepo:       repositories.AccountRepository,
			OperationTypeRepo: repositories.OperationRepository,
			TransctionRepo:    repositories.TransactionRepository,
		}),
	}
}

func handlersV1(c *config.AppConfig) v1.ApplicationHandler {
	repositoires := applicationRepositories(c)
	services := applicationServices(repositoires)
	return v1.ApplicationHandler{
		AccountHandler:     v1.NewAccountHandler(services.AccountService),
		TransactionHandler: v1.NewTransactionHandler(services.TransactionService),
	}
}

func main() {
	c := config.NewAppConfig()

	router.GET("/healthcheck", handler.HealthCheck)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			handerlsV1 := handlersV1(c)
			v1.POST("/accounts", handerlsV1.AccountHandler.Create)
			v1.GET("/accounts/:accountId", handerlsV1.AccountHandler.GetAccount)
			v1.POST("/transactions", handerlsV1.TransactionHandler.Create)
		}
	}

	router.Run(c.ListnerAddr())
}
