package handlers

import (
	"awesomeProject1/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

func DeleteCarPage(c echo.Context) error {
	kenteken := c.Param("kenteken")
	if kenteken == "" {
		data := echo.Map{"error": "No kenteken provided"}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	car, err := database.GetCarByKenteken(kenteken)
	if err != nil {
		data := echo.Map{"error": err}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	return c.Render(http.StatusOK, "deletecar.jet.html", car)
}

func DeleteCar(c echo.Context) error {
	kenteken := c.Param("kenteken")
	if kenteken == "" {
		data := echo.Map{"error": "No kenteken provided"}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	err := database.DeleteCar(kenteken)
	if err != nil {
		data := echo.Map{"error": err}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
