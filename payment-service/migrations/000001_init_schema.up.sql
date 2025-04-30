-- Create Mock Payments Table (Stores card and phone details)
CREATE TABLE mock_payments (
    id SERIAL PRIMARY KEY,
    card_holder TEXT,
    card_number TEXT UNIQUE,
    phone_number TEXT UNIQUE,
    expiry_date TEXT, -- Made nullable to support EWallet entries
    cvv TEXT,
    balance DECIMAL(10,2) DEFAULT 0 CHECK (balance >= 0),
    pin TEXT, -- for EWallet validation
    is_valid BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Payments Table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    student_id TEXT NOT NULL,
    payment_method TEXT NOT NULL CHECK (payment_method IN ('Card', 'EWallet')), -- Simplified payment methods
    total_amount DECIMAL(10,2) NOT NULL CHECK (total_amount > 0),
    transaction_id TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed')), -- Payment status
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    error_message TEXT -- error_message for failed payments
);

-- Create Payment Items Table (Stores individual items for each payment)
CREATE TABLE payment_items (
    id SERIAL PRIMARY KEY,
    payment_id INT NOT NULL REFERENCES payments(id) ON DELETE CASCADE, -- Foreign key to payments table
    course_id TEXT NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price > 0)
);

-- Create a function to update the updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Add triggers to update the updated_at column
CREATE TRIGGER set_mock_payments_updated_at
BEFORE UPDATE ON mock_payments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER set_payments_updated_at
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
