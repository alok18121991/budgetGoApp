package routes

import (
	"alok/web-service-budget/controllers"

	"github.com/labstack/echo/v4"
)

func ExpenseRoute(e *echo.Echo) {
	e.POST("/expense", controllers.CreateExpense)
	e.GET("/expense/:id", controllers.GetExpense)
	e.GET("/expense/:userId", controllers.GetAllExpenseForUser)
	e.DELETE("/expense/:id", controllers.DeleteExpense)
	e.DELETE("expense/all", controllers.DeleteAllExpense)

}
