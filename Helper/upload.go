package Helper

import (
	"github.com/labstack/echo"
	"io"
	"os"
	"path/filepath"
)

func FileUpload(c echo.Context) (string, error) {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("foto")
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Destination
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fileLocation := filepath.Join(dir, "files", file.Filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	// Copy
	if _, err = io.Copy(targetFile, src); err != nil {
		return "", err
	}
	return file.Filename, nil
}
