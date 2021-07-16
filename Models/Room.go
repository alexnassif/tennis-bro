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

func GetRoomsByUsers(user1_id string, privateRoom *[]Room)(err error) {
	
	if err = Config.DB.Where("user1_id = ?", user1_id).Or("user2_id = ?", user1_id).Find(&privateRoom).Error; err != nil {
		return err
	}
	return nil
}