package main

import (
	"project/config"
	"project/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	config.ConnectDatabase()
	r := routes.SetupRouter()
	r.Run(":8080")
}
