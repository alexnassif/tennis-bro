package main

import (
	"fmt"

	"github.com/alexnassif/tennis-bro/Config"
	"github.com/alexnassif/tennis-bro/Models"
	"github.com/alexnassif/tennis-bro/Routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error

func main() {
	dsn := "host=localhost user=alex password=1234 dbname=tennis_bro port=5432"
	Config.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer Config.Close()
	Config.DB.AutoMigrate(&Models.User{})
	r := Routes.SetupRouter()
	//running
	r.Run()
}
