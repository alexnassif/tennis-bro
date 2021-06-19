package Models

import (
	"strconv"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string      `json:"user_name"`
	Email    string      `json:"email"`
	Location Location    `gorm:"embedded" json:"location"`
	Level    PlayerLevel `gorm:"embedded" json:"level"`
	Image    string      `json:"image"`
	Bio      string      `json:"bio"`
	Password string      `json:"password"`
}

type PlayerLevel struct {
	Utr  float32
	Ntrp float32
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (user *User) GetId() string {
	return strconv.FormatUint(uint64(user.ID), 10)
}

func (user *User) GetName() string {
	return user.UserName
}

type OnlineUser interface {
	GetId() string
	GetName() string
}

type OnlineClient struct {
	ID       string `json:"id"`
	UserName string `json:"name"`
}

func (user *OnlineClient) GetId() string {
	return user.ID
}

func (user *OnlineClient) GetName() string {
	return user.UserName
}
