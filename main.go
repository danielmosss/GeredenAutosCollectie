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
)

func main() {
	godotenv.Load()

	e := echo.New()
	t := &Template{
		View: jet.NewSet(jet.NewOSFileSystemLoader("./views"), jet.InDevelopmentMode()),
	}
	e.Renderer = t

	// Set up routes
	e.GET("/", func(c echo.Context) error {
		cars, err := database.GetAllCars()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get cars")
		}

		return c.Render(http.StatusOK, "index.jet.html", cars)

		//v, err := views.GetTemplate("index.jet.html")
		//if err != nil {
		//	return c.String(http.StatusInternalServerError, "Failed to load template")
		//}
		//
		//var buf bytes.Buffer
		//if err := v.Execute(&buf, nil, cars); err != nil {
		//	return c.String(http.StatusInternalServerError, "Failed to execute template")
		//}
		//
		//return c.HTML(http.StatusOK, buf.String())
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

	e.GET("test-error", func(c echo.Context) error {
		err := c.Render(http.StatusOK, "error.jet.html", nil)
		if err != nil {
			c.Logger().Error(err)
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return nil
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
