-- +goose Up
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY, 
    amount FLOAT NOT NULL,
    account_id INT NOT NULL,
    operation_type_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(account_id) REFERENCES accounts(id),
    FOREIGN KEY(operation_type_id) REFERENCES operation_types(id)
);


-- +goose Down
DROP TABLE transactions;
