package kitchen

import (
	"kitchen/src/components/types/apparatus"
	"kitchen/src/components/types/cook"
	"kitchen/src/components/types/food"
	"kitchen/src/components/types/order"
)

type Kitchen struct {
	Cooks []cook.Cook
	Apparatus []apparatus.Apparatus
	Menu []food.Food
	OrderMap map[int16]order.Order
}
