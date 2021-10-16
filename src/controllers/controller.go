package controllers

import (
	coreService "kitchen/src/services"

	"github.com/gin-gonic/gin"
)

func processOrder(c *gin.Context) {
	// items := c.Query("items")
	// priority := c.Query("priority")
	// maxWait := c.Query("maxWait")

	id, err := c.Cookie("id")
	if(err != nil) {}

	c.JSON(200, gin.H{
		"id": id,
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

	ginEngine.GET("/order", processOrder)
}
