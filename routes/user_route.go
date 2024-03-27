package routes

import (
	"alok/web-service-budget/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	e.GET("/", controllers.GetAlbums)
	e.POST("/user", controllers.CreateUser)
	e.GET("/user/:userId", controllers.GetUser)
	e.GET("/user/email/:emailId", controllers.GetUserByEmail)
	e.GET("/user/all", controllers.GetAllUsers)
	e.DELETE("/user/all", controllers.DeleteAllUser)

}
