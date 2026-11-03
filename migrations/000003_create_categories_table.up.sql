-- +migrate Up
CREATE TABLE
  categories (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE, -- nullable untuk kategori global
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL, -- income, expense
    icon VARCHAR(50),
    created_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP
    WITH
      TIME ZONE DEFAULT CURRENT_TIMESTAMP
  );

CREATE INDEX idx_categories_user_id ON categories (user_id);