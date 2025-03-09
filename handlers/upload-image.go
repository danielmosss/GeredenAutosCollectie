package handlers

import (
	"awesomeProject1/database"
	"bytes"
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func UploadImage(c echo.Context) error {
	kenteken := c.FormValue("kenteken")

	// Retrieve file from form
	file, err := c.FormFile("picture")
	if err != nil {
		return c.String(http.StatusBadRequest, "Failed to get file")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to open file")
	}
	defer src.Close()

	// Read file content into a buffer
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, src); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read file")
	}

	// Encode file content to Base64
	encodedImage := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Call function to store the Base64 string in the database
	if err := database.AddPictureToCar(kenteken, encodedImage); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to add picture")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}
