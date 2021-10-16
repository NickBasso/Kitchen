package main

import (
	"kitchen/src/controllers"
	coreService "kitchen/src/services"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	router := gin.Default()

	coreService.InitCoreService()
	controllers.SetupController(router)

	router.Run(":4006")
}
