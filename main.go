package main

import (
	"awesomeProject1/database"
	"awesomeProject1/handlers"
	"bytes"
	"github.com/CloudyKit/jet/v6"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	godotenv.Load()
	// Initialize Jet views
	views := jet.NewSet(
		jet.NewOSFileSystemLoader("./views"),
		jet.InDevelopmentMode(),
	)

	e := echo.New()

	// Set up routes
	e.GET("/", func(c echo.Context) error {
		cars, err := database.GetAllCars()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get cars")
		}

		v, err := views.GetTemplate("index.jet")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to load template")
		}

		var buf bytes.Buffer
		if err := v.Execute(&buf, nil, cars); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to execute template")
		}

		return c.HTML(http.StatusOK, buf.String())
	})

	e.POST("/add-car", func(c echo.Context) error {
		return handlers.AddCar(c)
	})

	e.POST("/upload-image", func(c echo.Context) error {
		return handlers.UploadImage(c)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
