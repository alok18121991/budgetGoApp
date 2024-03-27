package services

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/models"
	"alok/web-service-budget/responses"
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserByEmail(c echo.Context, emailId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	err := configs.UserCollection.FindOne(ctx, bson.M{"email": emailId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func handleResponse(c echo.Context, map_ *echo.Map, message string, status int) error {
	return c.JSON(status, responses.GenericResponse{Status: status, Message: message, Data: map_})
}
