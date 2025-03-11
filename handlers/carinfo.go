package handlers

import (
	"awesomeProject1/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Carinfo(c echo.Context) error {
	kenteken := c.Param("kenteken")
	if kenteken == "" {
		data := echo.Map{
			"error": "No kenteken provided",
		}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	// Get car info from database
	car, err := database.GetCarByKenteken(kenteken)
	if err != nil {
		data := echo.Map{
			"error": err,
		}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	return c.Render(http.StatusOK, "carinfo.jet.html", car)
}
