-- +goose Up
-- Create Accounts Table
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance DECIMAL(10, 2) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Transactions Table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_id INT NOT NULL REFERENCES accounts(id),
    to_id INT NOT NULL REFERENCES accounts(id),
    amount DECIMAL(10, 2) NOT NULL CHECK (amount > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for Performance
CREATE INDEX idx_transactions_from_id ON transactions(from_id);
CREATE INDEX idx_transactions_to_id ON transactions(to_id);

-- +goose Down
-- Drop Indexes
DROP INDEX IF EXISTS idx_transactions_to_id;
DROP INDEX IF EXISTS idx_transactions_from_id;

-- Drop Tables
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS accounts;