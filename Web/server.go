package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Serve the index.html file
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
