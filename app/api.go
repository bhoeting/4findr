package main

import (
	"net/http"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func (app *App) run() {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Get("/professor", func(c *echo.Context) error {
		var professors []Professor
		app.DB.Find(&professors)
		return c.JSON(http.StatusOK, professors)
	})

	e.Get("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, 10)
	})

	e.Run(":3000")
}
