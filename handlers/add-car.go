package handlers

import (
	"awesomeProject1/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddCar(c echo.Context) error {
	kenteken := c.FormValue("kenteken")

	car, err := database.GetCarDataFromRDWAPI(kenteken)
	if err != nil {
		data := echo.Map{
			"error": err,
		}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	if err := database.SaveCarDrivenData(car); err != nil {
		data := echo.Map{
			"error": err,
		}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
