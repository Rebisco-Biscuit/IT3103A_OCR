package resolver

import (
	"log"
	"payment-mod/payment-service/graph/model"
	"sync"
)

var (
	Mu                     sync.Mutex
	CartUpdatedChannels    = make(map[string]chan []*model.CartItem)
	PaymentCreatedChannels = make(map[string]chan *model.Payment)
)

func BroadcastCartUpdate(studentId string, cart []*model.CartItem) {
	Mu.Lock()
	defer Mu.Unlock()

	if ch, exists := CartUpdatedChannels[studentId]; exists {
		log.Println("Broadcasting cart update for", studentId)
		ch <- cart
	} else {
		log.Println("No active subscription for", studentId)
	}
}
