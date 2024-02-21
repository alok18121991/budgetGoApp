package main

import (
	"alok/web-service-budget/configs"
	"alok/web-service-budget/routes"
	"fmt"
	"net/http"

	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	APP_ENV := "APP_ENV"
	if getEnv(APP_ENV) == "prod" {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}))
	} else {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
		}))
	}

	routes.UserRoute(e)
	routes.ExpenseRoute(e)
	configs.ConnectMongoDB()

	fmt.Println("app env :", getEnv(APP_ENV))

	if getEnv(APP_ENV) == "prod" {
		if err := e.StartTLS(":8080", getEnv("APP_SERVER_CERT"), getEnv("APP_SERVER_KEY")); err != http.ErrServerClosed {
			e.Logger.Fatal(err)
		}
	} else {
		e.Logger.Fatal(e.Start(":8080"))
	}
}

func getEnv(property string) string {
	return os.Getenv(property)
}
