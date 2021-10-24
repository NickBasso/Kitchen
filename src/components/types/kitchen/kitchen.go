package kitchen

import (
	"kitchen/src/components/types/cook"
	"kitchen/src/components/types/food"
	"kitchen/src/components/types/order"
)

type Kitchen struct {
	Cooks []cook.Cook
	Apparatus map[string]int
	Menu map[int]food.Food
	OrderMap map[string]order.Order
}
