-- +migrate Down
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_type;