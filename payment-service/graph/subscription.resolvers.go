package graph

import (
	"context"
	"payment-mod/payment-service/graph/model"
)

func (r *subscriptionResolver) PaymentCreated(ctx context.Context, studentId string) (<-chan *model.Payment, error) {
	ch := make(chan *model.Payment, 1)

	r.Mu.Lock()
	if _, exists := r.PaymentCreatedChannels[studentId]; !exists {
		r.PaymentCreatedChannels[studentId] = []chan *model.Payment{}
	}
	r.PaymentCreatedChannels[studentId] = append(r.PaymentCreatedChannels[studentId], ch)
	r.Mu.Unlock()

	go func() {
		<-ctx.Done()
		r.Mu.Lock()
		for i, subscriber := range r.PaymentCreatedChannels[studentId] {
			if subscriber == ch {
				r.PaymentCreatedChannels[studentId] = append(r.PaymentCreatedChannels[studentId][:i], r.PaymentCreatedChannels[studentId][i+1:]...)
				break
			}
		}
		r.Mu.Unlock()
		close(ch)
	}()

	return ch, nil
}
