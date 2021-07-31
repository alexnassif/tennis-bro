package Routes

import (
	
	"github.com/gin-gonic/gin"
	"github.com/alexnassif/tennis-bro/Controllers"
	"github.com/gin-contrib/cors"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	userGroup := r.Group("/user-api")
	{
		userGroup.GET("user", Controllers.GetUsers)
		userGroup.POST("user", Controllers.CreateUser)
		userGroup.GET("user/:id", Controllers.GetUserByID)
		userGroup.PUT("user/:id", Controllers.UpdateUser)
		userGroup.DELETE("user/:id", Controllers.DeleteUser)
	}

	roomGroup := r.Group("/room-api")
	{
		roomGroup.GET("rooms/:id", Controllers.GetRoomsForUser)
	}

	messageGroup := r.Group("/message-api")
	{
		messageGroup.GET("messages/:sender/:recipient", Controllers.GetMessagesByUserId)
	}
	return r
}
