package resolver

import (
	"payment-mod/payment-service/graph/model"
	"sync"
)

var (
	Mu                  sync.Mutex
	CartUpdatedChannels = make(map[string]chan []*model.CartItem)
)
