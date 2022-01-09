package services

import (
	"fmt"
	"kitchen/src/components/types/apparatus"
	"kitchen/src/components/types/cook"
	"kitchen/src/components/types/food"
	"kitchen/src/components/types/kitchen"
	"kitchen/src/components/types/mutextedValue"
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
	MutextedValue = mutextedValue.MutextedValue
)

var (
	kitchenRef *Kitchen = nil
  foodMenu = food.GetMenuMap()
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
			Mutex: new(sync.Mutex),
		},
		{
			ID: 1,
			Rank: 2,
			Proficiency: 3,
			Name: "Steve Peterson",
			CatchPhrase: "That's what she said..",
			WorkingCount: 0,
			Mutex: new(sync.Mutex),
		},
		{
			ID: 2,
			Rank: 2,
			Proficiency: 2,
			Name: "John Smithson",
			CatchPhrase: "How did you even get here lol ;=>",
			WorkingCount: 0,
			Mutex: new(sync.Mutex),
		},
		{
			ID: 3,
			Rank: 1,
			Proficiency: 2,
			Name: "Peter Owler",
			CatchPhrase: "Who's that even!?",
			WorkingCount: 0,
			Mutex: new(sync.Mutex),
		},
	}
	
	kitchenRef.Menu = food.GetMenuMap()
	kitchenRef.Apparatus = make(map[string]MutextedValue, 3)
	kitchenRef.Apparatus["None"] = MutextedValue{Value: 999, Mutex: new(sync.Mutex)}
	kitchenRef.Apparatus["Stove"] = MutextedValue{Value: 2, Mutex: new(sync.Mutex)}
	kitchenRef.Apparatus["Oven"] = MutextedValue{Value: 2, Mutex: new(sync.Mutex)}

	kitchenRef.OrderMap = make(map[string]Order)
	
	println("Cooks:\n")
	for _, cook := range kitchenRef.Cooks {
		fmt.Printf("%v\n", cook)
	}
	println("Apparatus:")
	fmt.Printf("%v\n", kitchenRef.Apparatus)
}

func ProcessOrder(order Order) Delivery {	
	println("ProcessOrder entered!")

	itemsCnt := len(order.Items)

	// wait group to await for all go routines to finish
	var waitGroup sync.WaitGroup
	// delivery channel to gather a dishes cooking results
	deliveryChannel := make(chan DeliveryCookingDetail, itemsCnt)
	
	cookedItems := make([]DeliveryCookingDetail, itemsCnt)

	fmt.Printf("ORDER ITEMS : %v\n", order.Items)
	
	waitGroup.Add(itemsCnt)
	fmt.Printf("wait groups: %d", itemsCnt)
	for i := 0; i < itemsCnt; i++ {
		go cookItem(order.Items[i], deliveryChannel, &waitGroup)
	}
	waitGroup.Wait()

	for i := 0; i < itemsCnt; i++ {
		cookedItems[i] = <-deliveryChannel

		fmt.Printf("Cooked item apparatus: %s\n", string(foodMenu[cookedItems[i].FoodID].Apparatus))
		fmt.Println(cookedItems[i])
	}

	println("SHOULD BE PRINTED AFTER WHOLE ORDER IS READY!")
  println(time.Now().UnixMilli() - order.PickUpTime)
	return Delivery { 
		OrderID: order.OrderID, TableID: order.TableID, WaiterID: order.WaiterID,
		Items: order.Items, Priority: order.Priority, MaxWait: order.MaxWait, PickUpTime: order.PickUpTime,
		CookingTime: time.Now().UnixMilli() - order.PickUpTime, CookingDetails: cookedItems,
  }  
}

func cookItem (dishID int, deliveryChannel chan DeliveryCookingDetail, waitGroup *sync.WaitGroup) { 
	println("cookItem entered!")

	fmt.Printf("FOOD ID: %d\n", dishID)
	fmt.Printf("Food: %v\n", foodMenu[dishID])

	isDishReady := false
	dishComplexity := foodMenu[dishID].Complexity
	foodApparatus := foodMenu[dishID].Apparatus

	apparatus := kitchenRef.Apparatus[string(foodApparatus)]

	apparatus.Mutex.Lock()
	apparatusAvailable := kitchenRef.Apparatus[string(foodApparatus)]
	apparatus.Mutex.Unlock()
	
	println("kitchen cooks: ", len(kitchenRef.Cooks))
	for !isDishReady {
		println("In loop!")
		
		for k := 0; k < 4; k++ {
			fmt.Printf("%v\n", kitchenRef.Cooks[k])
		}

		for i := 0; !isDishReady && i < len(kitchenRef.Cooks); i++ {
			cook := &kitchenRef.Cooks[i]

			if( (cook.Rank == dishComplexity ||
					cook.Rank - 1 == dishComplexity) &&
					cook.Proficiency > cook.WorkingCount &&
					apparatusAvailable.Value > 0) {
				defer waitGroup.Done()
			  cook.WorkingCount++
						
				apparatus.Mutex.Lock()
				apparatus = 
					MutextedValue {Mutex: kitchenRef.Apparatus[string(foodApparatus)].Mutex, Value: kitchenRef.Apparatus[string(foodApparatus)].Value - 1}
				apparatus.Mutex.Unlock()
				
				time.Sleep(time.Duration(foodMenu[dishID].PreparationTime) * time.Second)
				
				cook.Mutex.Lock()
				cook.WorkingCount--
				cook.Mutex.Unlock()

				apparatus.Mutex.Lock()
				apparatus = 
					MutextedValue {Mutex: kitchenRef.Apparatus[string(foodApparatus)].Mutex, Value: kitchenRef.Apparatus[string(foodApparatus)].Value - 1}
				apparatus.Mutex.Unlock()
				
				deliveryChannel <- DeliveryCookingDetail{FoodID: dishID, CookID: cook.ID}
				isDishReady = true
			}
		}
		
		time.Sleep(500 * time.Millisecond)
	}	

	fmt.Printf("Dish %d is ready!\n", dishID)
}