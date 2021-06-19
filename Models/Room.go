package Models

import (
	"github.com/alexnassif/tennis-bro/Config"
)

func AddRoom(room *Room) (err error) {
	if err = Config.DB.Create(room).Error; err != nil {
		return err
	}
	return nil
}

func FindRoomByName(room *Room, name string) (err error) {
	if err = Config.DB.Where("name = ?", name).First(room).Error; err != nil {
		return err
	}
	return nil
}