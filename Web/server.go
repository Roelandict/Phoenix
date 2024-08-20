package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Serve de index.html file
	e.GET("/", func(c echo.Context) error {
		return c.File("website/index.html")
	})

	// Serve statische bestanden (zoals CSS, JS, afbeeldingen) in de /website map
	e.Static("/css", "website/css")

	e.Logger.Fatal(e.Start(":80"))
}
