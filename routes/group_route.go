package routes

import (
	"alok/web-service-budget/controllers"

	"github.com/labstack/echo/v4"
)

func GroupRoute(e *echo.Echo) {
	e.POST("/group", controllers.CreateGroupHandler)
	e.POST("/group/add-owners", controllers.AddOwnersToGroupHandler)
	e.POST("/group/mark-inactive", controllers.MarkGroupInactiveHandler)
	e.GET("/group/:id", controllers.GetGroupDetailsHandler)
	e.GET("/groups", controllers.GetActiveGroupDetailsHandler)
}
