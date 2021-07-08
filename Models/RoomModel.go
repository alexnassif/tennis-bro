package Models

import "gorm.io/gorm"

type Room struct {
	gorm.Model
	Name    string
	Private bool
	User1ID int 
	User1 User `gorm:"foreignKey:User1ID"`
	User2ID int
	User2 User `gorm:"foreignKey:User2ID"`
}

func (room *Room) GetId() uint {
	return room.ID
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}
