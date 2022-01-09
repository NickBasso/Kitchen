package mutextedOrderMap

import (
	"kitchen/src/components/types/order"
	"sync"
)

type MutextedOrderMap struct {
	OrderMap 	map[string]order.Order
	Mutex 	  *sync.Mutex
}
