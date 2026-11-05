-- +migrate Up
CREATE TABLE budgets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    category_id BIGSERIAL NOT NULL,
    amount BIGINT NOT NULL,
    period VARCHAR(7) NOT NULL, -- Format: YYYY-MM
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_category_period UNIQUE(user_id, category_id, period)
);

CREATE INDEX idx_budgets_user_period ON budgets(user_id, period);