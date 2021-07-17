package Models

import (
	"strconv"

	"github.com/alexnassif/tennis-bro/Config"
	"gorm.io/gorm/clause"
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

func GetRoomsByUsers(user1_id string, privateRoom *[]Room) (err error) {
	u, _ := strconv.ParseUint(user1_id, 10, 32)
	if err = Config.DB.Preload(clause.Associations).Where("user1_id = ?", u).Or("user2_id = ?", u).Find(&privateRoom).Error; err != nil {
		return err
	}
	return nil
}
