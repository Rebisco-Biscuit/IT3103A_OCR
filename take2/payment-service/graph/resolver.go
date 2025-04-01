package graph

import (
	"context"
	"database/sql"
	"fmt"
	"payment-mod/payment-service/graph/model"
	"time"
)

type Resolver struct {
	DB *sql.DB
}

// GetPayment queries a specific payment by ID
func (r *queryResolver) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	// Query the database for a payment
	row := r.Resolver.DB.QueryRow("SELECT id, student_id, amount, currency, transaction_id, status, created_at FROM payments WHERE id = $1", id)

	var payment model.Payment
	err := row.Scan(&payment.ID, &payment.StudentID, &payment.Amount, &payment.Currency, &payment.TransactionID, &payment.Status, &payment.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("could not find payment with id %s: %v", id, err)
	}
	return &payment, nil
}

// ListPayments fetches all payments
func (r *queryResolver) ListPayments(ctx context.Context) ([]*model.Payment, error) {
	rows, err := r.Resolver.DB.Query("SELECT id, student_id, amount, currency, transaction_id, status, created_at FROM payments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.ID, &payment.StudentID, &payment.Amount, &payment.Currency, &payment.TransactionID, &payment.Status, &payment.CreatedAt); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}
	return payments, nil
}

// CreatePayment creates a new payment record
func (r *mutationResolver) CreatePayment(ctx context.Context, studentId string, amount float64, currency string, paymentMethod string, cardHolder *string, cardNumber *string, phoneNumber *string) (*model.Payment, error) {
	// Insert payment into the database
	var paymentID string
	createdAt := time.Now().Format(time.RFC3339)
	transactionID := fmt.Sprintf("TXN-%d", time.Now().Unix()) // Simple transaction ID (could be more complex)
	status := "pending"

	// Insert payment into the payments table
	err := r.Resolver.DB.QueryRow(`
		INSERT INTO payments (student_id, amount, currency, transaction_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		studentId, amount, currency, transactionID, status, createdAt).Scan(&paymentID)

	if err != nil {
		return nil, fmt.Errorf("could not create payment: %v", err)
	}

	// Insert into mock_payments table based on payment method
	if paymentMethod == "card" && cardHolder != nil && cardNumber != nil {
		_, err := r.Resolver.DB.Exec(`
			INSERT INTO mock_payments (payment_method, card_holder, card_number, phone_number, expiry_date, cvv)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			paymentMethod, *cardHolder, *cardNumber, "", "", "") // Expiry date and CVV left empty for simplicity
		if err != nil {
			return nil, fmt.Errorf("could not create mock payment for card: %v", err)
		}
	} else if paymentMethod == "gcash" && phoneNumber != nil {
		_, err := r.Resolver.DB.Exec(`
			INSERT INTO mock_payments (payment_method, card_holder, card_number, phone_number, expiry_date, cvv)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			paymentMethod, "", "", *phoneNumber, "", "") // Insert phone number for gcash payment method
		if err != nil {
			return nil, fmt.Errorf("could not create mock payment for gcash: %v", err)
		}
	}

	// Return created Payment object
	return &model.Payment{
		ID:            paymentID,
		StudentID:     studentId,
		Amount:        amount,
		Currency:      currency,
		TransactionID: transactionID,
		Status:        status,
		CreatedAt:     createdAt,
	}, nil
}
