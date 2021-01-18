package Models

import ("gorm.io/gorm")

type User struct {
	gorm.Model
	UserName string
	Latitute float64
	Longitude float64
	Email string
	Level string
	Image string
	Bio string
}

