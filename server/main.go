package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
)

// Index handler
func hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, world!\n")
	}
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Get("/", hello())

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	e.Run(standard.New(":" + port))
}
