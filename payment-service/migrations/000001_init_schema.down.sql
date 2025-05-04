-- Down migration: Drop triggers and tables

-- Drop triggers
DROP TRIGGER IF EXISTS set_mock_payments_updated_at ON mock_payments;
DROP TRIGGER IF EXISTS set_payments_updated_at ON payments;

-- Drop function for updating the updated_at column
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop tables
DROP TABLE IF EXISTS payment_items;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS mock_payments;


-- migrate -database "postgres://postgres:admin@localhost:5432/payment-service?sslmode=disable" -path migrations up