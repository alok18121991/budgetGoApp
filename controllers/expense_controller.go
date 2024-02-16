package controllers

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/models"
	"alok/web-service-budget/responses"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	limitString := c.Param("limit")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		limit = 5
	}
	defer cancel()

	var optionsParam *options.FindOptions
	if limit > 0 {
		optionsParam = options.Find().SetSort(bson.D{{Key: "expenseDate", Value: -1}}).SetLimit(int64(limit))
	} else {
		optionsParam = options.Find().SetSort(bson.D{{Key: "expenseDate", Value: -1}})
	}

	results, err := expenseCollection.Find(ctx, bson.M{"user_id": userId}, optionsParam)

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

func GetExpenseGroupByDate(c echo.Context) error {
	userIDsString := c.QueryParam("userids")
	numMonthsString := c.QueryParam("numMonths")
	groupType := c.QueryParam("groupType")

	if groupType == "" {
		return handleResponse(c, &echo.Map{"data": nil}, "userIds or groupType are empty", http.StatusBadRequest)
	}

	if groupType == "date" {
		groupType = "expenseDate"
	} else if groupType == "mode" {
		groupType = "source"
	} else {
		groupType = "type"
	}

	// groupField := "$" + groupType
	// Check if userIDsString is empty
	// Convert numMonthsString to an integer
	numMonths, shouldReturn, returnValue := validateRequest(userIDsString, c, numMonthsString)
	if shouldReturn {
		return returnValue
	}

	// Split user IDs string by comma
	userIDs := strings.Split(userIDsString, ",")
	for i, userID := range userIDs {
		userIDs[i] = strings.TrimSpace(userID)
	}

	currentDate := time.Now()
	currentMonth := currentDate.Month()
	currentYear := currentDate.Year()

	// Define the start and end of the current month
	startOfCurrentMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, time.UTC)
	// Get the month from the current date
	startOfMonth := startOfCurrentMonth.AddDate(0, -numMonths+1, 0) // Subtract numMonths-1 months
	startOfMonth = time.Date(startOfMonth.Year(), startOfMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := currentDate.AddDate(0, 1, 0).Add(-time.Nanosecond) // End of current month

	filter := bson.M{
		"user_id": bson.M{"$in": userIDs},
		"expenseDate": bson.M{
			"$gte": startOfMonth,
			"$lt":  endOfMonth,
		},
	}

	// Define the group field based on the groupType parameter
	var groupField interface{}
	switch groupType {
	case "expenseDate":
		groupField = bson.D{{Key: "$dateToString", Value: bson.D{
			{Key: "format", Value: "%Y-%m-%d"},
			{Key: "date", Value: "$expenseDate"},
		}}}
	default:
		groupField = "$" + groupType
	}

	// MongoDB aggregation pipeline to group by date and sum amount
	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: filter},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: groupField},
				{Key: "totalAmount", Value: bson.D{
					{Key: "$sum", Value: "$amount"},
				}},
			}},
		},
	}

	// Execute the aggregation pipeline
	cursor, err := expenseCollection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusInternalServerError)
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and store results
	results := make(map[string]float64)

	for cursor.Next(context.Background()) {
		var result struct {
			Type        string  `bson:"_id"`
			TotalAmount float64 `bson:"totalAmount"`
		}
		if err := cursor.Decode(&result); err != nil {
			return handleResponse(c, &echo.Map{"data": err.Error()}, "error", http.StatusInternalServerError)
		}
		results[result.Type] = result.TotalAmount
	}

	return handleResponse(c, &echo.Map{"data": results}, "success", http.StatusOK)
}

func validateRequest(userIDsString string, c echo.Context, numMonthsString string) (int, bool, error) {
	if userIDsString == "" {
		return 0, true, handleResponse(c, &echo.Map{"data": nil}, "userIds are empty", http.StatusInternalServerError)
	}

	numMonths, err := strconv.Atoi(numMonthsString)
	if err != nil || numMonths < 1 {
		return 0, true, handleResponse(c, &echo.Map{"data": nil}, "invalid numMonths", http.StatusBadRequest)
	}
	return numMonths, false, nil
}
