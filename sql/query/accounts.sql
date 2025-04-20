-- name: CreateAccount :one
INSERT INTO accounts(number) values ($1)
RETURNING *;

-- name: GetAccountById :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetAccountByNumber :one
SELECT * FROM accounts WHERE number = $1;
