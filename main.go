package main

import (
	"awesomeProject1/database"
	"awesomeProject1/handlers"
	"github.com/CloudyKit/jet/v6"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"strconv"
)

func main() {
	godotenv.Load()

	e := echo.New()
	t := &Template{
		View: jet.NewSet(jet.NewOSFileSystemLoader("./views"), jet.InDevelopmentMode()),
	}
	e.Renderer = t
	e.Debug = true

	// Set up routes
	e.GET("/", func(c echo.Context) error {
		cars, err := database.GetAllCars()
		if err != nil {
			data := echo.Map{
				"error": err,
			}
			return c.Render(http.StatusInternalServerError, "error.jet.html", data)
		}

		AmountOfCarsWithImages := 0
		for _, car := range cars {
			if car.Picture != "" {
				AmountOfCarsWithImages++
			}
		}

		data := echo.Map{
			"Cars":        cars,
			"TotalCars":   len(cars),
			"TotalImages": AmountOfCarsWithImages,
		}

		return c.Render(http.StatusOK, "index.jet.html", data)
	})

	e.GET("/car/:kenteken", func(c echo.Context) error {
		return handlers.Carinfo(c)
	})

	e.POST("/add-car", func(c echo.Context) error {
		return handlers.AddCar(c)
	})

	e.POST("/upload-image", func(c echo.Context) error {
		return handlers.UploadImage(c)
	})

	e.POST("/delete-image", func(c echo.Context) error {
		return handlers.DeleteImage(c)
	})

	e.GET("/edit-car/:kenteken", func(c echo.Context) error {
		kenteken := c.Param("kenteken")
		car, err := database.GetCarByKenteken(kenteken)
		if err != nil {
			data := echo.Map{
				"error": err,
			}
			return c.Render(http.StatusInternalServerError, "error.jet.html", data)
		}
		return c.Render(http.StatusOK, "editcar.jet.html", car)
	})

	e.POST("/update-car/:kenteken", func(c echo.Context) error {
		kenteken := c.Param("kenteken")

		var updatedCar, _ = database.GetCarByKenteken(kenteken)
		updatedCar.Kenteken = kenteken
		updatedCar.Merk = c.FormValue("merk")
		updatedCar.Handelsbenaming = c.FormValue("handelsbenaming")
		updatedCar.Variant = c.FormValue("variant")
		updatedCar.Uitvoering = c.FormValue("uitvoering")
		updatedCar.EersteKleur = c.FormValue("kleur")

		updatedCar.AantalZitplaatsen, _ = strconv.Atoi(c.FormValue("zitplaatsen"))
		updatedCar.AantalDeuren, _ = strconv.Atoi(c.FormValue("deuren"))
		updatedCar.AantalCilinders, _ = strconv.Atoi(c.FormValue("cilinders"))
		updatedCar.Catalogusprijs, _ = strconv.Atoi(c.FormValue("catalogusprijs"))

		err := database.UpdateCarData(kenteken, updatedCar)
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, "/")
	})

	e.Logger.Fatal(e.Start(":8080"))
}

type Template struct {
	View *jet.Set
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	template, err := t.View.GetTemplate(name)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	vars := make(jet.VarMap)
	if c.Get("flash") != nil {
		vars.Set("flash", c.Get("flash"))
	}
	err = template.Execute(w, vars, data)
	if err != nil {
		return err
	}
	return nil
}
