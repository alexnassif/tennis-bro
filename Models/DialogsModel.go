package Models

type Dialogs struct{
	User1ID int 
	User1 User `gorm:"foreignKey:User1ID"`
	User2ID int 
	User2 User `gorm:"foreignKey:User2ID"`
}