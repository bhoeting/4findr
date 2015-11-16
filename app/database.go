package main

import (
	"fmt"
	"log"
	"time"
	"unicode"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

// Subject is a category courses fall into.
// ex: CSE, MTH, ENG.
type Subject struct {
	ID     int    `json:"_id"`
	Title  string `json:"title"`
	Short  string `json:"short"`
	DeptID uint   `json:"dept_id"`

	DeletedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// Course is a type of class offered.
// ex: CSE174, MTH151, ENG111.
type Course struct {
	ID        int    `json:"_id"`
	Title     string `json:"title"`
	ShortName string `json:"short_name"`

	Subject   Subject `json:"-"`
	SubjectID int     `json:"subject_id"`

	DeletedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// Professor is the instructor
// of a course.
type Professor struct {
	ID   int    `json:"_id"`
	Name string `json:"name"`

	DeletedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// Class is a course that
// occured and has an
// average GPA.
type Class struct {
	ID  int     `json:"_id"`
	GPA float64 `json:"gpa" gorm:"column:gpa"`

	Course   Course `json:"-"`
	CourseID int    `json:"course_id"`

	Professor   Professor `json:"-"`
	ProfessorID int       `json:"professor_id"`

	DeletedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}

// ProfCoursePair is a struct that combines
// each class taught by the same professor
// and averages the GPA
type ProfCoursePair struct {
	ID  int     `json:"_id"`
	GPA float64 `json:"gpa" gorm:"column:gpa"`

	Course   Course `json:"course"`
	CourseID int    `json:"course_id"`

	Professor   Professor `json:"professor"`
	ProfessorID int       `json:"professor_id"`

	DeletedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
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

	if !app.DB.HasTable(&ProfCoursePair{}) {
		app.DB.CreateTable(&ProfCoursePair{})
	}

	return nil
}

func (app *App) findProfCoursePairsOrderedByGPA(courses []string) []ProfCoursePair {
	// TODO: Trash the fucking ORM.
	var profCoursePairs []ProfCoursePair
	rows, err := app.DB.Raw(
		`SELECT
		pairs.id, pairs.gpa, pairs.course_id, pairs.professor_id,
		courses.id, courses.title, courses.short_name,
		professors.id, professors.name
		FROM prof_course_pairs as pairs
		INNER JOIN professors ON professors.id = pairs.professor_id
		INNER JOIN courses ON courses.id = pairs.course_id
		AND courses.short_name IN (?)
		ORDER BY gpa DESC`, courses).Rows()

	if err != nil {
		log.Fatal(err.Error())
	}

	defer rows.Close()
	for rows.Next() {
		var profCoursePair ProfCoursePair
		rows.Scan(
			&profCoursePair.ID,
			&profCoursePair.GPA,
			&profCoursePair.CourseID,
			&profCoursePair.ProfessorID,
			&profCoursePair.Course.ID,
			&profCoursePair.Course.Title,
			&profCoursePair.Course.ShortName,
			&profCoursePair.Professor.ID,
			&profCoursePair.Professor.Name)

		profCoursePairs = append(profCoursePairs, profCoursePair)
	}

	return profCoursePairs
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
		shortName := fmt.Sprintf("%s%d", classData.ShortName, classData.Number)
		var course Course
		app.DB.Where(Course{ShortName: shortName}).
			Attrs(Course{Title: classData.Title, ShortName: shortName, SubjectID: subject.ID}).
			FirstOrCreate(&course)

		// Create the professor if they do not exist.
		var professor Professor
		app.DB.Where(Professor{Name: classData.Professor}).
			Attrs(Professor{Name: classData.Professor}).
			FirstOrCreate(&professor)

		// Create the class if it does not exist.
		var class Class
		app.DB.FirstOrCreate(&class,
			Class{CourseID: course.ID, ProfessorID: professor.ID, GPA: classData.GPA})

		// Create the professor-course pairs.
		var profCoursePair ProfCoursePair
		app.DB.FirstOrCreate(&profCoursePair, ProfCoursePair{CourseID: course.ID, ProfessorID: professor.ID})
	}

	// Set the GPA of all the professor-course-pairs.
	var profCoursePairs []ProfCoursePair
	app.DB.Find(&profCoursePairs)
	for _, profCoursePair := range profCoursePairs {
		type Result struct {
			Avg float64
		}

		var averageGPA Result
		app.DB.Raw("SELECT AVG(gpa) as avg FROM classes WHERE course_id = ? AND professor_id = ?",
			profCoursePair.CourseID, profCoursePair.ProfessorID).Scan(&averageGPA)

		app.DB.Model(&profCoursePair).Update("gpa", averageGPA.Avg)
	}

	return nil
}
