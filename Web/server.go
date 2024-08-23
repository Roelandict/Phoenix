package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Homepage
	e.GET("/homepage", func(c echo.Context) error {
		return c.File("website/index.html")
	})

	Redirect(e)
	//
	// Serve statische bestanden
	e.Static("/css", "website/css")

	e.Logger.Fatal(e.Start(":80"))
}
