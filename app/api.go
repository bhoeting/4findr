package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func getCollectionJSON(i interface{}, c *echo.Context) error {
	type CollectionJson struct {
		Count   int         `json:"count"`
		Results interface{} `json:"results"`
	}

	count := -1
	switch reflect.ValueOf(i).Kind() {
	case reflect.Slice:
		count = reflect.ValueOf(i).Len()
	default:
		// TODO: handle this in a useful way
		panic("Not slice")
	}

	b, _ := json.MarshalIndent(&CollectionJson{
		Count:   count,
		Results: i,
	}, "", "    ")

	c.Response().Header().Set(echo.ContentType, echo.ApplicationJSONCharsetUTF8)
	c.Response().Write(b)

	return nil
}

func (app *App) run() {
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	// api/v1/professor-coures-pairs?courses=ENG111,MTH151
	e.Get("/api/v1/professor-course-pairs", func(c *echo.Context) error {
		courses := strings.Split(c.Request().URL.Query().Get("courses"), ",")
		return getCollectionJSON(app.findProfCoursePairsOrderedByGPA(courses), c)
	})

	e.Get("/api/v1/professors", func(c *echo.Context) error {
		var professors []Professor
		app.DB.Find(&professors)
		return getCollectionJSON(professors, c)
	})

	e.Get("/api/v1/courses", func(c *echo.Context) error {
		var courses []Course
		app.DB.Find(&courses)
		return getCollectionJSON(courses, c)
	})

	e.Get("/api/v1/subjects", func(c *echo.Context) error {
		var subjects []Subject
		app.DB.Find(&subjects)
		return getCollectionJSON(subjects, c)
	})

	e.Get("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, 10)
	})

	e.Run(":3000")
}
