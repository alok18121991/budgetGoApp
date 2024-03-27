package controllers

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/models"
	"alok/web-service-budget/responses"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// var configs.UserCollection *mongo.Collection = configs.GetCollection(configs.DB, "user")
var validate = validator.New()

func CreateUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	//validate the request body
	if err := c.Bind(&user); err != nil {
		fmt.Println("Bind failed for object")
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusBadRequest)
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		fmt.Println("Valdiation failed for object")
		return handleResponse(c, &echo.Map{"data": validationErr.Error()}, "error", http.StatusBadRequest)
	}

	newUser := models.SetNewUserId(&user)

	result, err := configs.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusBadRequest)
	}

	return handleResponse(c, &echo.Map{"data": result}, "success", http.StatusOK)
}

func GetUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Param("userId")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	err := configs.UserCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusInternalServerError)
	}
	return handleResponse(c, &echo.Map{"data": user}, "success", http.StatusOK)

}

func GetUserByEmail(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	email := c.Param("emailId")
	var user models.User
	defer cancel()

	err := configs.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusInternalServerError)
	}
	return handleResponse(c, &echo.Map{"data": user}, "success", http.StatusOK)

}

func GetAllUsers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	defer cancel()

	results, err := configs.UserCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser models.User
		if err = results.Decode(&singleUser); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	return c.JSON(http.StatusOK, responses.GenericResponse{Status: http.StatusOK, Message: "success", Data: &echo.Map{"data": users}})
}

func DeleteAllUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteResult, err := configs.UserCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusInternalServerError)
	}
	return handleResponse(c, &echo.Map{"data": deleteResult}, "success", http.StatusOK)
}

func handleResponse(c echo.Context, map_ *echo.Map, message string, status int) error {
	return c.JSON(status, responses.GenericResponse{Status: status, Message: message, Data: map_})
}
