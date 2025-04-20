-- +goose Up
CREATE TABLE operation_types (
    id SERIAL PRIMARY KEY, 
    description TEXT NOT NULL UNIQUE,
    mode TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down 
DROP TABLE operation_types;
