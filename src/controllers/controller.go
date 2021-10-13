package controllers

import (
	coreService "kitchen/src/services"

	"github.com/gin-gonic/gin"
)

func placeOrder(c *gin.Context) {
	id := c.Query("id")
	items := c.Query("items")
	priority := c.Query("priority")
	maxWait := c.Query("maxWait")

	c.JSON(200, gin.H{
		"id":       id,
		"items":    items,
		"priority": priority,
		"maxWait":  maxWait,
		"MenuLength": len(coreService.GetMenu()),
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
		c.JSON(200, "Dining hall server is up!")
	})

	ginEngine.POST("/order", placeOrder)
}
