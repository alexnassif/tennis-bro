package main

import (
	"fmt"

	"github.com/alexnassif/tennis-bro/Auth"
	"github.com/alexnassif/tennis-bro/Config"
	"github.com/alexnassif/tennis-bro/Models"
	"github.com/alexnassif/tennis-bro/Routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var err error

func main() {
	dsn := "host=localhost dbname=tennis_bro port=5432"
	Config.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Status:", err)
	}
	defer Config.Close()
	Config.DB.AutoMigrate(&Models.User{}, &Models.OnlineClient{}, &Models.Room{}, &Models.Message{},)
	r := Routes.SetupRouter()
	//running
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	Config.CreateRedisClient()
	r.LoadHTMLFiles("index.html")

	r.GET("/room/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", Auth.AuthMiddleware(func(c *gin.Context) {

		ServeWs(wsServer, c.Writer, c.Request)

	}))

	r.GET("/api/login", func(c *gin.Context) {
		HandleLogin(c.Writer, c.Request)
	})
	
	r.Run()
}
