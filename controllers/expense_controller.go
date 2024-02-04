package controllers

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/models"
	"alok/web-service-budget/responses"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var expenseCollection *mongo.Collection = configs.GetCollection(configs.DB, "expenses")

func CreateExpense(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var expense models.Expense
	defer cancel()

	//validate the request body
	if err := c.Bind(&expense); err != nil {
		fmt.Println("Bind failed for object")
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusBadRequest)
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&expense); validationErr != nil {
		fmt.Println("Valdiation failed for object")
		return handleResponse(c, &echo.Map{"data": validationErr.Error()}, "error", http.StatusBadRequest)
	}

	newExpense := models.SetNewExpenseId(&expense)

	result, err := expenseCollection.InsertOne(ctx, newExpense)
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusBadRequest)
	}

	return handleResponse(c, &echo.Map{"data": result}, "success", http.StatusOK)
}

func GetExpense(c echo.Context) error {
	return nil
}

func GetAllExpenseForUser(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var expenses []models.Expense
	userId := c.Param("userId")
	defer cancel()

	results, err := expenseCollection.Find(ctx, bson.M{"user_id": userId})

	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusBadRequest)
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var expense models.Expense
		if err = results.Decode(&expense); err != nil {
			return c.JSON(http.StatusInternalServerError, responses.GenericResponse{Status: http.StatusInternalServerError, Message: "error", Data: &echo.Map{"data": err.Error()}})
		}

		expenses = append(expenses, expense)
	}

	return handleResponse(c, &echo.Map{"data": expenses}, "success", http.StatusOK)
}

func DeleteAllExpense(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteResult, err := expenseCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusInternalServerError)
	}
	return handleResponse(c, &echo.Map{"data": deleteResult}, "success", http.StatusOK)
}

func DeleteExpense(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id := c.Param("id")
	defer cancel()

	expenseId, _ := primitive.ObjectIDFromHex(id)
	deleteResult, err := expenseCollection.DeleteOne(ctx, bson.M{"_id": expenseId})
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusInternalServerError)
	}
	return handleResponse(c, &echo.Map{"data": deleteResult}, "success", http.StatusOK)
}
