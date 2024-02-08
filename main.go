package main

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/routes"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	routes.UserRoute(e)
	routes.ExpenseRoute(e)
	configs.ConnectMongoDB()

	e.Logger.Fatal(e.Start(":8080"))
}
