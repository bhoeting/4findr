package main

import "github.com/jinzhu/gorm"

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
