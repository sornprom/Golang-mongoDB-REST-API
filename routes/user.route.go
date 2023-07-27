package routes

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine) {
	//All routes related to
	userCollection := configs.GetCollection(configs.DB, "user")
	userController := controllers.NewUser(userCollection)
	router.POST("/user", userController.CreateUser)
	router.GET("/user/:userId", userController.GetAUser)
	router.PUT("/user/:userId", userController.EditAUser)
	router.DELETE("/user/:userId", userController.DeleteAUser)
	router.GET("/users", userController.GetAllUsers)
}
