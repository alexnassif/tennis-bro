package main
import (
 "github.com/alexnassif/tennis-bro/Config"
 "github.com/alexnassif/tennis-bro/Models"
 "github.com/alexnassif/tennis-bro/Routes"
 "fmt"

"gorm.io/driver/postgres"
"gorm.io/gorm"
)
var err error
func main() {
	dsn := "host=localhost user=alexnassif password=123 dbname=tennis_bro port=5432"
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