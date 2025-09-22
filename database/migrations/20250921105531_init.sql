-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(255),
    age SMALLINT,
    country VARCHAR(2)
);

-- +goose Down
DROP TABLE IF EXISTS users;