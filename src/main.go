package main

import (
	"kitchen/src/controllers"
	"kitchen/src/services"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.ForceConsoleColor()
	router := gin.Default()
  
	services.InitCoreService()
	controllers.SetupController(router)

	router.Run(":4006")
}
