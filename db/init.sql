CREATE TABLE IF NOT EXISTS transaction (
    id SERIAL PRIMARY KEY,
    transaction_hash VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    block_number NUMERIC NOT NULL,
    data JSON NOT NULL
);