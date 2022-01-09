package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kitchen/src/components/types/order"
	"kitchen/src/services"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type (
	Order = order.Order
  Delivery = order.Delivery
)

func processOrder(c *gin.Context) {
	var order order.Order;

	jsonDataRaw, _ := ioutil.ReadAll(c.Request.Body)
	_ = json.Unmarshal(jsonDataRaw, &order)

	fmt.Printf("POST order %s received, processing...\n", order.OrderID)
	c.JSON(200, "Kitchen: Order received, processing...");

	delivery := services.ProcessOrder(order)

	reqBody, reqBodySerializationErr := json.Marshal(delivery)
		if reqBodySerializationErr != nil {
			log.Fatalln(reqBodySerializationErr)
		}

	println(os.Getenv("DHALL_URL")+"/distribution");
		
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

func SetupController(ginEngine *gin.Engine) {
	ginEngine.GET("/", func(c *gin.Context) {
		c.JSON(200, "Kitchen server is up!")
	})

	ginEngine.POST("/order", processOrder)
}
