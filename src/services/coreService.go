package services

import (
	"fmt"
	"kitchen/src/components/constants"
	"kitchen/src/components/types/apparatus"
	"kitchen/src/components/types/cook"
	"kitchen/src/components/types/food"
	"kitchen/src/components/types/kitchen"
	"kitchen/src/components/types/order"
	"kitchen/src/components/types/order/props"
)

type Kitchen = kitchen.Kitchen
type Cook = cook.Cook
type Apparatus = apparatus.Apparatus
type ApparatusState = apparatus.ApparatusState
type Food = food.Food
type Order = order.Order
type ItemCookingDetail = props.ItemCookingDetail
type Delivery = order.Delivery

var kitchenRef *Kitchen = nil
var foodMenu = food.GetMenuMap()

func InitCoreService() {
	kitchenRef = new(Kitchen)
	kitchenRef.Cooks = []Cook{
		{
			ID: 0,
			Rank: 3,
			Proficiency: 4,
			Name: "Gordon Ramsay",
			CatchPhrase: "Hey, panini head, are you listening to me?",
			WorkingCount: 0,
		},
		{
			ID: 1,
			Rank: 2,
			Proficiency: 2,
			Name: "Steve Peterson",
			CatchPhrase: "That's what she said..",
			WorkingCount: 0,
		},
		{
			ID: 2,
			Rank: 1,
			Proficiency: 2,
			Name: "John Smithson",
			CatchPhrase: "Who's that even!?",
			WorkingCount: 0,
		},
	}

	kitchenRef.Menu = food.GetMenuMap()
	kitchenRef.Apparatus = make(map[string]int, 3)
	kitchenRef.Apparatus["None"] = 999
	kitchenRef.Apparatus["Stove"] = 2
	kitchenRef.Apparatus["Oven"] = 1

	/* ApparatusState{
		Ovens: apparatus.OvenState{TotalCount: 1, FreeCount: 1},
		Stoves: apparatus.StoveState{TotalCount: 2, FreeCount: 2},
	} */
	kitchenRef.OrderMap = make(map[string]Order, constants.GeneratedOrdersCount)
	
	println("Cooks:\n")
	for _, cook := range kitchenRef.Cooks {
		fmt.Printf("%v\n", cook)
	}
	println("Apparatus:")
	fmt.Printf("%v\n", kitchenRef.Apparatus)
}

func ProcessOrder(order Order) Delivery {
	kitchenRef.OrderMap[order.OrderID] = order
	itemsCnt := len(order.Items)

	// orderChannel := make(chan Delivery)
	cookedItems := make([]ItemCookingDetail, itemsCnt)
	itemDetailChannels := make([]chan ItemCookingDetail, itemsCnt)
	for i := range itemDetailChannels {
		itemDetailChannels[i] = make(chan ItemCookingDetail)
 }

 println("channels size: ", len(itemDetailChannels))

	for i := 0; i < itemsCnt; i++ {
		go cookItem(order.Items[i], itemDetailChannels[i])
	}

	// go cookOrder(order, orderChannel)

	/* for k, v := range kitchenRef.OrderMap {
		
	} */
fmt.Printf("orderChannel ProcessOrder: %v\n", itemDetailChannels)
	for i := 0; i < itemsCnt; i++ {
		cookedItems[i] = <-itemDetailChannels[i]
		fmt.Printf("Channel %d result: %v\n", i, cookedItems[i])
		kitchenRef.Cooks[cookedItems[i].CookID].WorkingCount-- 
		kitchenRef.Apparatus[string(foodMenu[cookedItems[i].FoodID].Apparatus)]++
		//close(itemDetailChannels[i])
	}

	// close(itemDetailChannels)
	
	delete(kitchenRef.OrderMap, order.OrderID)
	return Delivery{OrderID: order.OrderID}
}

func cookItem(foodID int, itemCookingDetailChannel chan ItemCookingDetail)  {
	for i := 0; i < len(kitchenRef.Cooks); i++ {
		cook := kitchenRef.Cooks[i]
	  if( (cook.Rank == foodMenu[foodID].Complexity ||
				cook.Rank - 1 == foodMenu[foodID].Complexity) &&
				cook.Proficiency > cook.WorkingCount &&
				kitchenRef.Apparatus[string(foodMenu[foodID].Apparatus)] > 0){
					cook.WorkingCount++
					kitchenRef.Apparatus[string(foodMenu[foodID].Apparatus)]--

					itemCookingDetailChannel <- ItemCookingDetail{FoodID: foodID, CookID: cook.ID, CookingTime: kitchenRef.Menu[foodID].PreparationTime}
		}
	}

	// orderChannel <- Delivery{}
	// fmt.Printf("orderChannel cookOrder: %v\n", orderChannel)

	// return Delivery{}
}