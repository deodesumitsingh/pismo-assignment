-- name: CreateTransaction :one
INSERT INTO transactions(account_id, operation_type_id, amount) VALUES($1, $2, $3)
RETURNING *;
