package routes

import (
	"alok/web-service-budget/controllers"

	"github.com/labstack/echo/v4"
)

func AuthRRoute(e *echo.Echo) {
	e.POST("/auth/refresh", controllers.RefreshTokenHandler)
	e.GET("/auth/user", controllers.GetAuthUserDetails)
}
