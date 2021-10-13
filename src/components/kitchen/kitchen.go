package kitchen

import (
	"kitchen/src/components/apparatus"
	"kitchen/src/components/cook"
	"kitchen/src/components/food"
	"kitchen/src/components/order"
)

type Kitchen struct {
	Cooks []cook.Cook
	Apparatus []apparatus.Apparatus
	Menu []food.Food
	OrderMap map[int16]order.Order
}
