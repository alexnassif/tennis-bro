package Config

import(
	"gorm.io/gorm"
	"fmt"
)

var DB *gorm.DB

type DBConfig struct {
	Host string
	Port int
	User string
	DBName string
	Password string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host: "localhost",
		Port: 5432,
		User: "root",
		Password: "",
		DBName: "tennis-bro",
	}
	return &dbConfig
}

func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
	 "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
	 dbConfig.User,
	 dbConfig.Password,
	 dbConfig.Host,
	 dbConfig.Port,
	 dbConfig.DBName,
	)
   }

func Close() {
	tbDB, _ := DB.DB()
	tbDB.Close()
}