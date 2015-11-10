package main

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

// Subject is a category courses fall into.
// ex: CSE, MTH, ENG.
type Subject struct {
	ID     int
	Title  string
	Short  string
	DeptID uint

	DeletedAt time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Course is a type of class offered.
// ex: CSE174, MTH151, ENG111.
type Course struct {
	ID        int
	GPA       float64
	Title     string
	ShortName string

	Subject   Subject
	SubjectID sql.NullInt64

	DeletedAt time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Professor is the instructor
// of a course.
type Professor struct {
	ID   int
	Name string

	DeletedAt time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Class is a course that
// occured and has an
// average GPA.
type Class struct {
	ID int

	Course   Course
	CourseID sql.NullInt64

	Professor   Professor
	ProfessorID sql.NullInt64

	DeletedAt time.Time
	UpdatedAt time.Time
	CreatedAt time.Time
}

func (app *App) initDB() error {
	var err error
	if app.DB, err = gorm.Open("sqlite3", "/tmp/4findr.db"); err != nil {
		return err
	}

	if !app.DB.HasTable(&Subject{}) {
		app.DB.CreateTable(&Subject{})
	}

	if !app.DB.HasTable(&Course{}) {
		app.DB.CreateTable(&Course{})
	}

	if !app.DB.HasTable(&Professor{}) {
		app.DB.CreateTable(&Professor{})
	}

	if !app.DB.HasTable(&Class{}) {
		app.DB.CreateTable(&Class{})
	}

	return nil
}
