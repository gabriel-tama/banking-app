CREATE TABLE IF NOT EXISTS "balances"(
    id SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    balance INT NOT NULL,
    currencyCode CHAR(3) NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id),
    CONSTRAINT unique_balance UNIQUE (userId, currencyCode)
);

