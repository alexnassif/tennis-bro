package Models

type Room struct {
	Id      string
	Name    string
	Private bool
	User1ID int 
	User1 User `gorm:"foreignKey:User1ID"`
	User2ID int
	User2 User `gorm:"foreignKey:User2ID"`
}

func (room *Room) GetId() string {
	return room.Id
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) GetPrivate() bool {
	return room.Private
}
