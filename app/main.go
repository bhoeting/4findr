package main

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

//
// 2015
// Brennan Hoeting
// Oxford, OH
//

// App is a struct that stores all the
// components of the application
type App struct {
	DB gorm.DB
}

func main() {
	var app App

	if err := app.initDB(); err != nil {
		log.Fatal(err)
	}

	args := os.Args
	if len(args) > 1 {
		switch args[1] {
		case "fetch":
			subjects, classes := fetch()
			app.seedDB(subjects, classes)
		default:
			app.run()
		}
	} else {
		app.run()
	}
}
