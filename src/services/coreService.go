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
)

 type ( 
	Kitchen = kitchen.Kitchen
	Cook = cook.Cook
	Apparatus = apparatus.Apparatus
	ApparatusState = apparatus.ApparatusState
	Food = food.Food
	Order = order.Order
	ItemCookingDetail = props.ItemCookingDetail
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

	kitchenRef.OrderMap = make(map[string]Order, constants.GeneratedOrdersCount)
	
	println("Cooks:\n")
	for _, cook := range kitchenRef.Cooks {
		fmt.Printf("%v\n", cook)
	}
	println("Apparatus:")
	fmt.Printf("%v\n", kitchenRef.Apparatus)
}

func ProcessOrder(order Order) Delivery {
	println("ProcessOrder entered!");

	kitchenRef.OrderMap[order.OrderID] = order
	itemsCnt := len(order.Items)

	cookedItems := make([]ItemCookingDetail, itemsCnt)
	jobs := make(chan int, itemsCnt)
	results := make(chan ItemCookingDetail, itemsCnt)

	go cookItem(jobs, results)
	for i := 0; i < itemsCnt; i++ {
		jobs <- order.Items[i]
	}
	close(jobs)

	for j := 0; j < itemsCnt; j++ {
		cookedItems[j] = <-results
		fmt.Printf("Cooked item apparatus: %s\n", string(foodMenu[cookedItems[j].FoodID].Apparatus))
		kitchenRef.Cooks[cookedItems[j].CookID].WorkingCount = kitchenRef.Cooks[cookedItems[j].CookID].WorkingCount - 1
		println("ovens before: ", kitchenRef.Apparatus["Oven"])
		kitchenRef.Apparatus[string(foodMenu[cookedItems[j].FoodID].Apparatus)] = kitchenRef.Apparatus[string(foodMenu[cookedItems[j].FoodID].Apparatus)] + 1 
		println("ovens after: ", kitchenRef.Apparatus["Oven"])

		delete(kitchenRef.OrderMap, order.OrderID)
		fmt.Println(cookedItems[j])
	}

	defer foodMenuMutex.Unlock() 
	defer foodApparatusMutex.Unlock()
	defer apparatusMapMutex.Unlock()

	return Delivery{}
} 

func cookItem(foods <-chan int, results chan<- ItemCookingDetail) {
	println("cookItem entered!")
	println("foods channel size: ", len(foods))

	for foodID := range foods {
		println("inside range foods!")
		foodMenuMutex.Lock()
		defer foodMenuMutex.Unlock() 
		dishComplexity := foodMenu[foodID].Complexity

		foodApparatusMutex.Lock()
		defer foodApparatusMutex.Unlock()
		foodApparatus := foodMenu[foodID].Apparatus

		apparatusMapMutex.Lock()
		defer apparatusMapMutex.Unlock()
		apparatusAvailable := kitchenRef.Apparatus[string(foodApparatus)]
		
		println("kitchen cooks: ", len(kitchenRef.Cooks))
		for {
			for i := 0; i < len(kitchenRef.Cooks); i++ {
				cook := &kitchenRef.Cooks[i]

				// fmt.Printf("DishComplexity=%d - FoodApparatus=%s - ApparatusAvailable=%d - CookRank=%d - CookProficiency=%d \n",
				//  dishComplexity, foodApparatus, apparatusAvailable, cook.Rank, cook.Proficiency)

				// println("isTakingItem: ", 
												// (cook.Rank == dishComplexity ||
												// cook.Rank - 1 == dishComplexity) &&
												// cook.Proficiency > cook.WorkingCount &&
												// apparatusAvailable > 0)
				// fmt.Printf("%t %t %t\n", (cook.Rank == dishComplexity ||
					// cook.Rank - 1 == dishComplexity), cook.Proficiency > cook.WorkingCount, 	apparatusAvailable > 0)
				// println("apparatus used: ", string(foodApparatus))

				if( (cook.Rank == dishComplexity ||
						cook.Rank - 1 == dishComplexity) &&
						cook.Proficiency > cook.WorkingCount &&
						apparatusAvailable > 0){
							cook.WorkingCount++
							kitchenRef.Apparatus[string(foodApparatus)] = kitchenRef.Apparatus[string(foodApparatus)] - 1
							println("apparatus after taking item to cook: ", kitchenRef.Apparatus[string(foodApparatus)])
							
							results <- ItemCookingDetail{FoodID: foodID, CookID: cook.ID, CookingTime: kitchenRef.Menu[foodID].PreparationTime}
							break
				}
			}
		}
	}
}