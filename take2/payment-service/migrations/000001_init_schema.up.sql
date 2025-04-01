-- Create Payments Table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    student_id TEXT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency TEXT NOT NULL,
    transaction_id TEXT UNIQUE NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Mock Payments Table
CREATE TABLE mock_payments (
    id SERIAL PRIMARY KEY,
    payment_method TEXT NOT NULL,
    card_holder TEXT,
    card_number TEXT,
    phone_number TEXT,
    expiry_date TEXT,
    cvv TEXT
);
