# Pismo Assignment

This assignment is aim to provide **api's** for the transactions.  

## APIs Path 
`/healthcheck` - This path is responsible for checking the health of the application.  
`/api/v1/transactions/{id}` - This path is responsible for handling the transactions by id.  
`/api/v1/transactions` - This path is responsible for handling the transactions.   
`/api/v1/accounts` - This path is responsible for handling the accounts.   
`/api/v1/accounts/{id}` - This path is responsible for handling the accounts by id.  

## Implementation Details
This project is solely depend upon whether to use the `InMemory` for data repository or `Postgres` totally depending upon `DbURL` env variable.

If you are willing to use **PSQL** please ensure to have following steps - 
1. Postgres URL.
2. Install `goose` via `go install github.com/pressly/goose/v3/cmd/goose@latest` this will be use for running the migrtion.
3. All the SQL is generated with the use of **SQLC**. If you are willing to contribute or extend this project please install it via `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`

## Up and running the project
There are majorly 3 ENV variables 
- *HOST*: This will set the host for running the application. Defaults to `localhost`
- *PORT*: This will set the port at which application can run on. Defaults to `8080`
- *DbURL*: This is the database complete URL. Default to `""`. If you provide `DbURL` and application can't able to connect with it then application won't start.

> [!NOTE]
> If you are using **DbURL** before running the application please run the migration present inside sql/migrations using `goose`. **goose postgres `postgresql://<username>:<password>@<host>:<port>/<dbname> up`  

You can set these env variables using good old `export` command like `export HOST="127.0.0.1"` or you can create `.env` file in the root of the project and that will work fine too.  

Once, you have set all the data you can run the application for development as `make dev-run` or for PRODUCTION as `make run`.  

**Happy Coding**
