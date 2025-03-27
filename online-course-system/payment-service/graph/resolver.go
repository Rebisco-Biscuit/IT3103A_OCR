package graph

import (
	"context"
	"time"

	"payment-service/graph/model"
	"payment-service/prisma/db"
)

// CreatePayment resolver
func (r *mutationResolver) CreatePayment(ctx context.Context, studentId string, amount float64, currency string, transactionId string) (*model.Payment, error) {
	now := time.Now()

	payment, err := r.Resolver.Prisma.Payment.CreateOne(
		db.Payment.StudentID.Set(studentId),
		db.Payment.Amount.Set(amount),
		db.Payment.Currency.Set(currency),
		db.Payment.TransactionID.Set(transactionId),
		db.Payment.Status.Set("pending"),
		db.Payment.CreatedAt.Set(now),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	// Convert to GraphQL model
	return &model.Payment{
		ID:            payment.ID,
		StudentID:     payment.StudentID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
		Refund:        nil,
	}, nil
}

// GetPayment resolver
func (r *queryResolver) GetPayment(ctx context.Context, id string) (*model.Payment, error) {
	payment, err := r.Resolver.Prisma.Payment.FindUnique(
		db.Payment.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return &model.Payment{
		ID:            payment.ID,
		StudentID:     payment.StudentID,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		TransactionID: payment.TransactionID,
		Status:        payment.Status,
		CreatedAt:     payment.CreatedAt.Format(time.RFC3339),
		Refund:        nil,
	}, nil
}
