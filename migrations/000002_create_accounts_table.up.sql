-- +migrate Up
CREATE TABLE
  accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL, -- bank, cash, e-wallet
    balance BIGINT NOT NULL DEFAULT 0, -- disimpan dalam unit terkecil (misal: Rupiah)
    currency VARCHAR(3) NOT NULL DEFAULT 'IDR',
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );

CREATE INDEX idx_accounts_user_id ON accounts (user_id);