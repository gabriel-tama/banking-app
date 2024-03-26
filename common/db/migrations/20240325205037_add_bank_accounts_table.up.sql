CREATE TABLE IF NOT EXISTS "bankaccounts"(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    bank_name VARCHAR(255) NOT NULL ,
    account_number VARCHAR(255) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_bank_name ON bankaccounts (bank_name);
