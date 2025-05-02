package graph

import (
	"context"
	"database/sql"
	"fmt"
	"payment-mod/payment-service/graph/model"
	"sync"
	"time"
)

// Resolver serves as dependency injection for the app
type Resolver struct {
	DB *sql.DB // Database connection

	Mu                     sync.Mutex
	PaymentCreatedChannels map[string]chan *model.Payment
}

// Mutation Resolver
func (r *Resolver) CreatePayment(ctx context.Context, studentID string, items []*model.PaymentItemInput, paymentMethod model.PaymentMethod, cardHolder, cardNumber, expiryDate, cvv, phoneNumber, pin *string) (*model.Payment, error) {
	fmt.Println("CreatePayment called")

	totalAmount := 0.0
	for _, item := range items {
		totalAmount += item.Price
	}

	status := "failed"

	// Common check for duplicate course payments
	for _, item := range items {
		exists, err := r.checkCourseAlreadyPaid(ctx, studentID, item.CourseID)
		if err != nil {
			return nil, fmt.Errorf("failed to validate existing course payments: %v", err)
		}
		if exists {
			return r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, fmt.Sprintf("%s is already paid.", item.CourseID))
		}
	}

	// Validate payment method specific things
	switch paymentMethod {
	case model.PaymentMethodCard:
		if cardHolder == nil || cardNumber == nil || expiryDate == nil || cvv == nil {
			return r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, "cardHolder, cardNumber, expiryDate, and cvv are required for card payments")
		}

		mockID, balance, err := r.validateMockPayment(ctx, "card", *cardHolder, *cardNumber, *expiryDate, *cvv, "", "")
		if err != nil {
			return r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, err.Error())
		}

		if balance < totalAmount {
			return r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, "insufficient balance on the card")
		}

		err = r.deductMockBalance(ctx, mockID, totalAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct balance: %v", err)
		}

		status = "completed"

	case model.PaymentMethodEWallet:
		if phoneNumber == nil || pin == nil {
			return r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, "phoneNumber and pin are required for e-wallet payments")
		}

		mockID, balance, err := r.validateMockPayment(ctx, "ewallet", "", "", "", "", *phoneNumber, *pin)
		if err != nil {
			return r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, err.Error())
		}

		if balance < totalAmount {
			return r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, "insufficient balance in the e-wallet")
		}

		err = r.deductMockBalance(ctx, mockID, totalAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to deduct balance: %v", err)
		}

		status = "completed"

	default:
		return nil, fmt.Errorf("unsupported payment method")
	}

	payment, err := r.insertPaymentWithItems(ctx, studentID, items, totalAmount, paymentMethod, status, "")
	if err != nil {
		return nil, err
	}

	// Publish the payment to the subscription channel
	r.Mu.Lock()
	if ch, ok := r.PaymentCreatedChannels[studentID]; ok {
		ch <- payment
	}
	r.Mu.Unlock()

	return payment, nil
}

// 2025 april 28 update: check if the course is already paid
func (r *Resolver) checkCourseAlreadyPaid(ctx context.Context, studentID, courseID string) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1
			FROM payments p
			JOIN payment_items pi ON p.id = pi.payment_id
			WHERE p.student_id = $1
			AND pi.course_id = $2
			AND p.status = 'completed'
		)
	`, studentID, courseID).Scan(&exists)
	return exists, err
}

// 2025 april 28 update: payment validation
func (r *Resolver) validateMockPayment(ctx context.Context, method, cardHolder, cardNumber, expiryDate, cvv, phoneNumber, pin string) (int, float64, error) {
	var mockID int
	var balance float64
	var query string
	var args []interface{}

	if method == "card" {
		query = `
			SELECT id, balance FROM mock_payments 
			WHERE card_holder = $1 AND card_number = $2 AND expiry_date = $3 AND cvv = $4 AND is_valid = TRUE`
		args = []interface{}{cardHolder, cardNumber, expiryDate, cvv}
	} else if method == "ewallet" {
		query = `
			SELECT id, balance FROM mock_payments 
			WHERE phone_number = $1 AND pin = $2 AND is_valid = TRUE`
		args = []interface{}{phoneNumber, pin}
	} else {
		return 0, 0, fmt.Errorf("unsupported payment method")
	}

	err := r.DB.QueryRow(query, args...).Scan(&mockID, &balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("%s not found or invalid", method)
		}
		return 0, 0, fmt.Errorf("failed to validate %s: %v", method, err)
	}

	return mockID, balance, nil
}

func (r *Resolver) deductMockBalance(ctx context.Context, mockID int, amount float64) error {
	_, err := r.DB.Exec(`
		UPDATE mock_payments 
		SET balance = balance - $1 
		WHERE id = $2`, amount, mockID)
	return err
}

// Helper function to insert payment and items into the database
func (r *Resolver) insertPaymentWithItems(ctx context.Context, studentID string, items []*model.PaymentItemInput, totalAmount float64, paymentMethod model.PaymentMethod, status string, errorMessage string) (*model.Payment, error) {
	// Insert payment into the payments table
	var paymentID string
	transactionID := fmt.Sprintf("TXN-%d", time.Now().Unix())
	createdAt := time.Now().Format(time.RFC3339)

	err := r.DB.QueryRowContext(ctx, `
		INSERT INTO payments (student_id, payment_method, total_amount, transaction_id, status, created_at, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		studentID, paymentMethod, totalAmount, transactionID, status, createdAt, errorMessage).Scan(&paymentID)

	if err != nil {
		fmt.Println("Payment insertion failed:", err) // Debug log
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	fmt.Println("Payment created with ID:", paymentID) // Debug log

	// Insert items into the payment_items table
	for _, item := range items {
		_, err := r.DB.ExecContext(ctx, `
			INSERT INTO payment_items (payment_id, course_id, price)
			VALUES ($1, $2, $3)`,
			paymentID, item.CourseID, item.Price)

		if err != nil {
			fmt.Println("Payment item insertion failed:", err) // Debug log
			return nil, fmt.Errorf("failed to insert payment item: %v", err)
		}
	}

	// Convert items to []*model.PaymentItem
	var paymentItems []*model.PaymentItem
	for _, item := range items {
		paymentItems = append(paymentItems, &model.PaymentItem{
			CourseID: item.CourseID,
			Price:    item.Price,
		})
	}

	// Return the created payment
	return &model.Payment{
		ID:            paymentID,
		StudentID:     studentID,
		Items:         paymentItems,
		TotalAmount:   totalAmount,
		TransactionID: transactionID,
		PaymentMethod: paymentMethod,
		Status:        model.PaymentStatus(status),
		CreatedAt:     createdAt,
		ErrorMessage:  &errorMessage, // Include errorMessage in the response
	}, nil
}

// Query Resolver for getPayment
func (r *Resolver) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	var payment model.Payment
	err := r.DB.QueryRow(`
        SELECT id, student_id, total_amount,transaction_id, payment_method, status, created_at 
        FROM payments WHERE id = $1`, id).
		Scan(&payment.ID, &payment.StudentID, &payment.TotalAmount, &payment.TransactionID, &payment.PaymentMethod, &payment.Status, &payment.CreatedAt)

	if err != nil {
		fmt.Printf("Error fetching payment with ID %s: %v\n", id, err)
		return nil, fmt.Errorf("could not find payment with id %s", id)
	}

	// Fetch payment items
	rows, err := r.DB.Query(`
        SELECT course_id, price FROM payment_items WHERE payment_id = $1`, id)
	if err != nil {
		fmt.Println("Error fetching payment items:", err)
		return nil, err
	}
	defer rows.Close()

	var items []*model.PaymentItem
	for rows.Next() {
		var item model.PaymentItem
		if err := rows.Scan(&item.CourseID, &item.Price); err != nil {
			fmt.Println("Error scanning payment item:", err)
			return nil, err
		}
		items = append(items, &item)
	}
	payment.Items = items

	return &payment, nil
}

// Query Resolver for listPayments
func (r *Resolver) ListPayments(ctx context.Context) ([]*model.Payment, error) {
	rows, err := r.DB.Query(`
		SELECT id, student_id, total_amount, transaction_id, payment_method, status, created_at 
		FROM payments`)
	if err != nil {
		fmt.Println("Error fetching payments:", err)
		return nil, err
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		var payment model.Payment
		if err := rows.Scan(&payment.ID, &payment.StudentID, &payment.TotalAmount, &payment.TransactionID, &payment.PaymentMethod, &payment.Status, &payment.CreatedAt); err != nil {
			fmt.Println("Error scanning payment:", err)
			return nil, err
		}

		// Fetch payment items for each payment
		itemRows, err := r.DB.Query(`
            SELECT course_id, price FROM payment_items WHERE payment_id = $1`, payment.ID)
		if err != nil {
			fmt.Println("Error fetching payment items:", err)
			return nil, err
		}
		defer itemRows.Close()

		var items []*model.PaymentItem
		for itemRows.Next() {
			var item model.PaymentItem
			if err := itemRows.Scan(&item.CourseID, &item.Price); err != nil {
				fmt.Println("Error scanning payment item:", err)
				return nil, err
			}
			items = append(items, &item)
		}
		payment.Items = items

		payments = append(payments, &payment)
	}

	return payments, nil
}

// Query Resolver for listPaymentHistory
func (r *Resolver) ListPaymentHistory(ctx context.Context, studentID string) ([]*model.PaymentHistory, error) {
	rows, err := r.DB.Query(`
		SELECT p.transaction_id, p.status, p.payment_method, p.created_at, pi.course_id, pi.price
		FROM payments p
		JOIN payment_items pi ON p.id = pi.payment_id
		WHERE p.student_id = $1 AND p.status = 'completed'
		ORDER BY p.created_at ASC`, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment history: %w", err)
	}
	defer rows.Close()

	var history []*model.PaymentHistory
	for rows.Next() {
		var record model.PaymentHistory
		if err := rows.Scan(&record.TransactionID, &record.Status, &record.PaymentMethod, &record.CreatedAt, &record.CourseID, &record.Price); err != nil {
			return nil, fmt.Errorf("failed to scan payment history: %w", err)
		}
		history = append(history, &record)
	}
	return history, nil
}
