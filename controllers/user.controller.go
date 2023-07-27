package controllers

import (
	"context"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

type user struct {
	c *mongo.Collection
}

func NewUser(c *mongo.Collection) *user {
	return &user{c: c}
}

type createUser struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
	Title    string `json:"title" binding:"required"`
}

func (u *user) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    gin.H{"data": err.Error()},
		})
		return
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    gin.H{"data": validationErr.Error()},
		})
		return
	}
	newUser := createUser{
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}
	result, err := u.c.InsertOne(ctx, newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    gin.H{"data": err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    gin.H{"data": result},
	},
	)
}

func (u *user) GetAUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := u.c.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    gin.H{"data": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    gin.H{"data": user}})
}

func (u *user) EditAUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    gin.H{"data": err.Error()},
		})
		return
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    gin.H{"data": validationErr.Error()},
		})
		return
	}

	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}
	result, err := u.c.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    gin.H{"data": err.Error()}})
		return
	}

	//get updated user details
	updatedUser := models.User{}
	if result.MatchedCount == 1 {
		err := u.c.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    gin.H{"data": err.Error()}})
			return
		}
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    gin.H{"data": updatedUser}})
}

func (u *user) DeleteAUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := u.c.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    gin.H{"data": err.Error()}})
		return
	}

	if result.DeletedCount < 1 {
		c.JSON(http.StatusNotFound,
			responses.UserResponse{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    gin.H{"data": "User with specified ID not found!"}},
		)
		return
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    gin.H{"data": "User successfully deleted!"}},
	)
}

func (u *user) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := u.c.Find(ctx, bson.M{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    gin.H{"data": err.Error()}})
		return
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    gin.H{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	c.JSON(http.StatusOK,
		responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    gin.H{"data": users}},
	)
}
