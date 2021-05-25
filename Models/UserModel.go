package Models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string
	Email    string
	Location Location    `gorm:"embedded"`
	Level    PlayerLevel `gorm:"embedded"`
	Image    string
	Bio      string
	Password string
}

type PlayerLevel struct {
	Utr  float32
	Ntrp float32
}

type Location struct {
	Latitude  float64
	Longitude float64
}
