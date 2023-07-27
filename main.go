package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRouter(r) // add this

	r.Run("localhost:6000")
}
