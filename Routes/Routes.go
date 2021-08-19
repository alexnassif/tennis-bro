package Routes

import (
	"github.com/alexnassif/tennis-bro/Auth"
	"github.com/alexnassif/tennis-bro/Controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	userGroup := r.Group("/user-api")
	{
		userGroup.GET("user", Controllers.GetUsers)
		userGroup.POST("user", Controllers.CreateUser)
		userGroup.GET("user/:id", Controllers.GetUserByID)
		userGroup.PUT("user/:id", Auth.AuthMiddleware(func(c *gin.Context) {
			Controllers.UpdateUser(c)
		}))
		userGroup.DELETE("user/:id", Auth.AuthMiddleware(func(c *gin.Context) {
			Controllers.DeleteUser(c)
		}))
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
