package main

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/routes"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	routes.UserRoute(e)
	routes.ExpenseRoute(e)
	configs.ConnectMongoDB()

	e.Logger.Fatal(e.Start(":8080"))
}
