package main

import (
	"kitchen/src/controllers"
	coreService "kitchen/src/services"

	"github.com/gin-gonic/gin"
)

func main() {
	ginEngine := gin.Default()

	coreService.InitCoreService()
	controllers.SetupController(ginEngine)

	ginEngine.Run(":4006")
}
