-- name: GetOperationById :one
SELECT * FROM operation_types WHERE id = $1;
