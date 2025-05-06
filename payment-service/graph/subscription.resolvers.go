package graph

import (
	"context"
	"log"
	"payment-mod/payment-service/graph/model"
)

func (r *subscriptionResolver) PaymentCreated(ctx context.Context, studentId string) (<-chan *model.Payment, error) {
	log.Println("Subscription triggered for studentId:", studentId)
	ch := make(chan *model.Payment, 1)

	r.Mu.Lock()
	if _, exists := r.PaymentCreatedChannels[studentId]; !exists {
		r.PaymentCreatedChannels[studentId] = ch
	}
	r.Mu.Unlock()

	go func() {
		<-ctx.Done()
		r.Mu.Lock()
		if r.PaymentCreatedChannels[studentId] == ch {
			delete(r.PaymentCreatedChannels, studentId)
		}
		r.Mu.Unlock()
		close(ch)
	}()

	return ch, nil
}

func (r *subscriptionResolver) CartUpdated(ctx context.Context, studentId string) (<-chan []*model.CartItem, error) {
	log.Println("Subscription triggered for cart updates for studentId:", studentId)
	ch := make(chan []*model.CartItem, 1)

	r.Mu.Lock()
	if _, exists := r.CartUpdatedChannels[studentId]; !exists {
		r.CartUpdatedChannels[studentId] = ch
	}
	r.Mu.Unlock()
	go func() {
		<-ctx.Done()
		r.Mu.Lock()
		if r.CartUpdatedChannels[studentId] == ch {
			delete(r.CartUpdatedChannels, studentId)
		}
		r.Mu.Unlock()
		close(ch)
	}()

	return ch, nil
}
