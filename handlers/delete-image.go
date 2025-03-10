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
		data := echo.Map{
			"error": err,
		}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
