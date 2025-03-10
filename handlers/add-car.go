package handlers

import (
	"awesomeProject1/database"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AddCar(c echo.Context) error {
	kenteken := c.FormValue("kenteken")

	// Get car data from RDW API
	car, err := database.GetCarDataFromRDWAPI(kenteken)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to get car data: %v", err))
	}

	// Save the car data
	if err := database.SaveCarDrivenData(car); err != nil {
		data := echo.Map{
			"error": err,
		}
		return c.Render(http.StatusInternalServerError, "error.jet.html", data)
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
