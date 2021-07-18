package Models

import (
	"github.com/alexnassif/tennis-bro/Config"
	"gorm.io/gorm/clause"
)

func GetMessagesByUser(sender_id string, recipient_id string, messages *[]Message)(err error){

	if err = Config.DB.Preload(clause.Associations).Where("recipient_id = ? AND sender_id = ?", recipient_id, sender_id).Or("recipient_id = ? AND sender_id = ?", sender_id, recipient_id).Find(&messages).Error; err != nil{
		return err
	}
	return nil
}