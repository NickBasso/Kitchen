package main

import (
	"kitchen/src/configs"
	"kitchen/src/controllers"
	"kitchen/src/services"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	router := gin.Default()

	configs.SetupENV()
	services.InitCoreService()
	controllers.SetupController(router)

	router.Run(":4006")
}
