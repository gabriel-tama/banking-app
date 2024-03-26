CREATE TABLE IF NOT EXISTS "transactions"(
    id SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    amount INT NOT NULL,
    balanceId INT NOT NULL,
    bankAccountNumber VARCHAR(30) NOT NULL,
    bankName VARCHAR(30) NOT NULL,
    transferProofImg VARCHAR(100),
    currencyCode CHAR(3) NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userId) REFERENCES users(id)
);

