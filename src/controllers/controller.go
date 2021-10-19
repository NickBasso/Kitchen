package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kitchen/src/components/types/order"
	coreService "kitchen/src/services"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Order = order.Order
type Delivery order.Delivery

func processOrder(c *gin.Context) {
	var order order.Order;

	jsonDataRaw, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {}

	e := json.Unmarshal(jsonDataRaw, &order)
	if e != nil {}

	fmt.Printf("POST order %s received, processing...\n", order.OrderID)
	c.JSON(200, "Order received, processing...");

	delivery := coreService.ProcessOrder(order)

	reqBody, reqBodySerializationErr := json.Marshal(delivery)
		if reqBodySerializationErr != nil {
			log.Fatalln(reqBodySerializationErr)
		}

	resp, POSTErr := http.Post(os.Getenv("DHALL_URL")+"/distribution", "application/json", bytes.NewBuffer(reqBody))
		if POSTErr != nil {
			log.Fatalln(POSTErr)
		}

	defer resp.Body.Close()

	body, readPOSTResErr := ioutil.ReadAll(resp.Body)
	if readPOSTResErr != nil {
		log.Fatalln(readPOSTResErr)
	}

	var POSTDeliveryRes string;
	POSTResDeserializationErr := json.Unmarshal(body, &POSTDeliveryRes)
	if(POSTResDeserializationErr != nil) {
		log.Fatalln(POSTResDeserializationErr)
	}

	fmt.Printf("POST delivery: %s => %v\n", delivery.OrderID, POSTDeliveryRes)
}

func getOrderList(c *gin.Context) {
	id := c.Query("id")
	items := c.Query("items")
	priority := c.Query("priority")
	maxWait := c.Query("maxWait")

	c.JSON(200, gin.H{
		"id":       id,
		"items":    items,
		"priority": priority,	
		"maxWait":  maxWait,
	})
}

func SetupController(ginEngine *gin.Engine) {
  coreService.InitCoreService();

	ginEngine.GET("/", func(c *gin.Context) {
		c.JSON(200, "Kitchen server is up!")
	})

	ginEngine.POST("/order", processOrder)
}
