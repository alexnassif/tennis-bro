package Models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Location Location    `gorm:"embedded" json:"location"`
	Level    PlayerLevel `gorm:"embedded" json:"level"`
	Image    string `json:"image"`
	Bio      string `json:"bio"`
	Password string `json:"password"`
}

type PlayerLevel struct {
	Utr  float32
	Ntrp float32
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Profile interface {
	GetId() string
	GetName() string
}

func (user *User) GetId() string {
    return string(user.ID)
}

func (user *User) GetName() string {
    return user.UserName
}