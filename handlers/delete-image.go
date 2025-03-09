package handlers

import (
	"awesomeProject1/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteImage(c echo.Context) error {
	kenteken := c.FormValue("kenteken")
	err := database.DeletePicterOfCar(kenteken)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to delete image")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
