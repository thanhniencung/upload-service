package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"path/filepath"
	"net/http"
	"github.com/twinj/uuid"
	"os"
	"fmt"
)

func upload(c echo.Context) error {

	images := make([]string, 0)

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
					"message": err.Error(),
			})
	}
	files := form.File["files"]

	for _, file := range files {
		// Source
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusOK, echo.Map{
					"message": err.Error(),
			})
		}
		defer src.Close()

	fileName := uuid.NewV4();
	filePath := fmt.Sprintf("images/product/%s%s", fileName, filepath.Ext(file.Filename))

		// Destination
		dst, err := os.Create(filePath)
		if err != nil {
			return c.JSON(http.StatusOK, echo.Map{
					"message": err.Error(),
			})
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusOK, echo.Map{
					"message": err.Error(),
			})
		}

		url := fmt.Sprintf("http://localhost:3004/static/product/%s%s", fileName, filepath.Ext(file.Filename))
		images = append(images, url)
	}

	fmt.Println(echo.Map{
		"images": images,
	})

	fmt.Println()
	fmt.Println()

	return c.JSON(http.StatusOK, echo.Map{
		"images": images,
	})
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "images/")
	e.POST("/upload", upload)

	e.Logger.Fatal(e.Start(":3004"))
}
