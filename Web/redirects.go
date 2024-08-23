package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Redirect(e *echo.Echo) {
	e.GET("/home", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/homepage")
	})
}
