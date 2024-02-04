package controllers

import (
	"alok/web-service-budget/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAlbums(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, types.GetAlbums(), " ")
}
