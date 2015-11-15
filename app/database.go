package main

import (
	"fmt"
	"time"
	"unicode"

	"github.com/bhoeting/out"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
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
	SubjectID int

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
	ID  int
	GPA float64

	Course   Course
	CourseID int

	Professor   Professor
	ProfessorID int

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

func (app *App) findSubjectByShort(short string) Subject {
	var subject Subject
	app.DB.Model(&Subject{}).Where("short = ?", short).First(&subject)
	return subject
}

func (app *App) seedDB(subjects []Subject, allClassData []ClassData) error {
	// Create the subjects.
	if subjectsCount := len(subjects); subjectsCount > 0 {
		var subjectsRowCount int
		app.DB.Model(&Subject{}).Count(subjectsRowCount)

		if subjectsRowCount < subjectsCount {
			for _, subject := range subjects {
				app.DB.Create(&subject)
			}
		}
	}

	fmt.Println("%d", len(allClassData))

	// Create the classes/courses/professors.
	for _, classData := range allClassData {
		// Extract the subject code, (`ENG111` -> `ENG`).
		short := ""
		for _, r := range classData.ShortName {
			if !unicode.IsLetter(rune(r)) {
				break
			}
			short += string(rune(r))
		}

		// Get the subject
		subject := app.findSubjectByShort(short)

		// Create the course if it does not exist.
		var course Course
		app.DB.Where(Course{ShortName: classData.ShortName}).
			Attrs(Course{Title: classData.Title, ShortName: classData.ShortName, SubjectID: subject.ID}).
			FirstOrCreate(&course)

		out.Out(course.ID)

		// Create the professor if they do not exist.
		var professor Professor
		app.DB.Where(Professor{Name: classData.Professor}).
			Attrs(Professor{Name: classData.Professor}).
			FirstOrCreate(&professor)

		// Create the class if it does not exist.
		var class Class
		app.DB.FirstOrCreate(&class, Class{CourseID: course.ID, ProfessorID: professor.ID, GPA: classData.GPA})
	}

	return nil
}
