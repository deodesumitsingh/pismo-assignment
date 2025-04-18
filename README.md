# Pismo Assignment

This assignment is aim to provide **api's** for the transactions.  

## How to run
I am relying on *make* as the build tool.  
`make dev-run` will start the project and ready for API request.  
`make test` will run the tests.  

## APIs Path 
`/healthcheck` - This path is responsible for checking the health of the application.  
`/api/v1/transactions/{id}` - This path is responsible for handling the transactions by id.  
`/api/v1/transactions` - This path is responsible for handling the transactions.   
`/api/v1/accounts` - This path is responsible for handling the accounts.   
`/api/v1/accounts/{id}` - This path is responsible for handling the accounts by id.  

## Implementation Details
**InMemoryDB** is used as the database for this project.  
**PSQL** usage is still in progress 🚧.  
