package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Homepage
	e.GET("/homepage", func(c echo.Context) error {
		return c.File("website/index.html")
	})
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/homepage")
	})
	Redirect(e)
	//
	// Serve statische bestanden
	e.Static("/css", "website/css")

	e.Logger.Fatal(e.Start(":80"))
}
