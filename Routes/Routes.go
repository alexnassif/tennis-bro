package Routes

import (
	"github.com/alexnassif/tennis-bro/Auth"
	"github.com/alexnassif/tennis-bro/Controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Authorization"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	r.Use(cors.New(config))

	userGroup := r.Group("/user-api")
	{
		//userGroup.GET("user", Controllers.GetUsers)

		userGroup.GET("user", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.GetUsers(c)
		}))
		userGroup.POST("user", Controllers.CreateUser)
		userGroup.GET("user/:id", Controllers.GetUserByID)
		userGroup.PUT("user/:id", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.UpdateUser(c)
		}))
		userGroup.DELETE("user/:id", Auth.AuthMiddlewareRest(func(c *gin.Context) {
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
