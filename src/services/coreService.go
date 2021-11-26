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
	"sync"
	"time"
)

 type ( 
	Kitchen = kitchen.Kitchen
	Cook = cook.Cook
	Apparatus = apparatus.Apparatus
	ApparatusState = apparatus.ApparatusState
	Food = food.Food
	Order = order.Order
	ItemCookingDetail = props.ItemCookingDetail
	DeliveryCookingDetail = props.DeliveryCookingDetail
	Delivery = order.Delivery
)


var (
	kitchenRef *Kitchen = nil
  foodMenu = food.GetMenuMap()
  foodMenuMutex = sync.RWMutex{}
	apparatusMapMutex = sync.RWMutex{}
	foodApparatusMutex = sync.RWMutex{}
	orderMapMutex = sync.RWMutex{}	
)

func InitCoreService() {
	kitchenRef = new(Kitchen)
	kitchenRef.Cooks = []Cook{
		/* {
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
		}, */
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
			Proficiency: 3,
			Name: "Steve Peterson",
			CatchPhrase: "That's what she said..",
			WorkingCount: 0,
		},
		{
			ID: 2,
			Rank: 2,
			Proficiency: 2,
			Name: "John Smithson",
			CatchPhrase: "How did you even get here lol ;=>",
			WorkingCount: 0,
		},
		{
			ID: 3,
			Rank: 1,
			Proficiency: 2,
			Name: "Peter Owler",
			CatchPhrase: "Who's that even!?",
			WorkingCount: 0,
		},
	}

	kitchenRef.Menu = food.GetMenuMap()
	kitchenRef.Apparatus = make(map[string]int, 3)
	kitchenRef.Apparatus["None"] = 999
	kitchenRef.Apparatus["Stove"] = 2
	kitchenRef.Apparatus["Oven"] = 2

	kitchenRef.OrderMap = make(map[string]Order, constants.GeneratedOrdersCount)
	
	println("Cooks:\n")
	for _, cook := range kitchenRef.Cooks {
		fmt.Printf("%v\n", cook)
	}
	println("Apparatus:")
	fmt.Printf("%v\n", kitchenRef.Apparatus)
}

func ProcessOrder(order Order) Delivery {	
	println("ProcessOrder entered!")

	kitchenRef.OrderMap[order.OrderID] = order
	itemsCnt := len(order.Items)
	cookedItems := make([]DeliveryCookingDetail, itemsCnt)

	fmt.Printf("ORDER ITEMS : %v\n", order.Items)
	
	deliveries := cookOrder(order.Items)

	for j := 0; j < itemsCnt; j++ {
		// cookedItems[j] = <-results
		fmt.Printf("Cooked item apparatus: %s\n", string(foodMenu[cookedItems[j].FoodID].Apparatus))
		// kitchenRef.Cooks[cookedItems[j].CookID].WorkingCount = kitchenRef.Cooks[cookedItems[j].CookID].WorkingCount - 1
		println("ovens before: ", kitchenRef.Apparatus["Oven"])
		// kitchenRef.Apparatus[string(foodMenu[cookedItems[j].FoodID].Apparatus)] = kitchenRef.Apparatus[string(foodMenu[cookedItems[j].FoodID].Apparatus)] + 1 
		println("ovens after: ", kitchenRef.Apparatus["Oven"])

		delete(kitchenRef.OrderMap, order.OrderID)
		fmt.Println(cookedItems[j])
	}

	// defer foodMenuMutex.Unlock() 
	// defer foodApparatusMutex.Unlock()
	// defer apparatusMapMutex.Unlock()

	println("SHOULD BE PRINTED AFTER ALL COOING IS DONE!")
println(time.Now().UnixMilli() - order.PickUpTime)
	return Delivery{ OrderID: order.OrderID, TableID: order.TableID, WaiterID: order.WaiterID,
		               Items: order.Items, Priority: order.Priority, MaxWait: order.MaxWait, PickUpTime: order.PickUpTime,
		               CookingTime: time.Now().UnixMilli() - order.PickUpTime, CookingDetails: deliveries,
        }  
	}

func cookOrder(foods []int) []DeliveryCookingDetail{
	fmt.Printf("FOODS: %v\n", foods)
	
	println("cookItem entered!")
	println("foods channel size: ", len(foods))

	readyCounter := 0
	deliveries := make([]DeliveryCookingDetail, len(foods))

	for j := 0; j < len(foods); j++ {
		foodID := foods[j]
		println("inside range foods!")
		fmt.Printf("FOOD ID: %d\n", foodID)
		fmt.Printf("Food: %v\n", foodMenu[foodID])
		
		// foodMenuMutex.Lock()
		// defer foodMenuMutex.Unlock() 
		dishComplexity := foodMenu[foodID].Complexity

		// foodApparatusMutex.Lock()
		// defer foodApparatusMutex.Unlock()
		foodApparatus := foodMenu[foodID].Apparatus

		// apparatusMapMutex.Lock()
		// defer apparatusMapMutex.Unlock()
		apparatusAvailable := kitchenRef.Apparatus[string(foodApparatus)]
		
		println("kitchen cooks: ", len(kitchenRef.Cooks))
		for readyCounter < len(foods) {

			for i := 0; readyCounter < len(foods) && i < len(kitchenRef.Cooks); i++ {
				cook := &kitchenRef.Cooks[i]

				// fmt.Printf("Before:\n\tDishComplexity=%d - FoodApparatus=%s - ApparatusAvailable=%d - CookRank=%d - CookProficiency=%d - CookWorkingCount=%d \n",
				//  dishComplexity, foodApparatus, apparatusAvailable, cook.Rank, cook.Proficiency, cook.WorkingCount)

				if( (cook.Rank == dishComplexity ||
						cook.Rank - 1 == dishComplexity) &&
						cook.Proficiency > cook.WorkingCount &&
						apparatusAvailable > 0){
							cook.WorkingCount++
							
							// apparatusMapMutex.RLock()
							// defer apparatusMapMutex.RUnlock()
							kitchenRef.Apparatus[string(foodApparatus)] = kitchenRef.Apparatus[string(foodApparatus)] - 1
							// fmt.Printf("apparatus(%s) after taking item to cook: %d\n", foodApparatus, kitchenRef.Apparatus[string(foodApparatus)])
							time.Sleep(time.Duration(foodMenu[foodID].PreparationTime) * time.Second)
							
							cook.WorkingCount--
							kitchenRef.Apparatus[string(foodApparatus)] = kitchenRef.Apparatus[string(foodApparatus)] + 1
							// fmt.Printf("apparatus(%s) after returning cooked item: %d\n", foodApparatus, kitchenRef.Apparatus[string(foodApparatus)])
								
							deliveries[readyCounter] = DeliveryCookingDetail{FoodID: foodID, CookID: cook.ID}
							readyCounter++

							// results <- DeliveryCookingDetail{FoodID: foodID, CookID: cook.ID}
				}
			}
		}
	}

	return deliveries
}