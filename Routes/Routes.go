package Routes

import (
	"github.com/alexnassif/tennis-bro/Controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

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

