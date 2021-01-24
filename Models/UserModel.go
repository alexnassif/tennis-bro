package Models

import ("gorm.io/gorm")

type User struct {
	gorm.Model
	UserName string
	Latitude float64
	Longitude float64
	Email string
	Level float64
	Image string
	Bio string
	Password string
}

