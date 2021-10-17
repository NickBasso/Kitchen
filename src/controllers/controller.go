package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kitchen/src/components/constants"
	"kitchen/src/services"
	coreService "kitchen/src/services"

	"github.com/gin-gonic/gin"
)

type Order services.Order

func processOrder(c *gin.Context) {
	// items := c.Query("items")
	// priority := c.Query("priority")
	// maxWait := c.Query("maxWait")

	// id, err := c.Cookie("id")
	// if(err != nil) {}
	orders := make([]Order, constants.GeneratedOrdersCount)
	jsonDataRaw, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {}

	e := json.Unmarshal(jsonDataRaw, &orders)
	if e != nil {}

	for i := 0; i < len(orders); i++ {
		fmt.Printf("Order %d: %v\n", i + 1, orders[i])
	}

	c.JSON(200, gin.H{
		"id": 412,
		"orders": c.Request.Body,
	})
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
	// init operations' service
  coreService.InitCoreService();

	// home path
	ginEngine.GET("/", func(c *gin.Context) {
		c.JSON(200, "Kitchen server is up!")
	})

	ginEngine.POST("/order", processOrder)
}
