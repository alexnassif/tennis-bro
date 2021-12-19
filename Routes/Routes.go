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
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}

	r.Use(cors.New(config))

	userGroup := r.Group("/user-api")
	{

		userGroup.GET("user", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.GetUsers(c)
		}))
		userGroup.POST("user", Controllers.CreateUser)
		userGroup.GET("user/:id", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.GetUserByID(c)
		}))
		userGroup.PUT("user/:id", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.UpdateUser(c)
		}))
		userGroup.DELETE("user/:id", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.DeleteUser(c)
		}))
		userGroup.GET("me", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.FetchUser(c)
		}))
	}

	roomGroup := r.Group("/room-api")
	{
		roomGroup.GET("rooms", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.GetRoomsForUser(c)
		}))
	}

	messageGroup := r.Group("/message-api")
	{
		messageGroup.GET("messages/:recipient", Auth.AuthMiddlewareRest(func(c *gin.Context) {
			Controllers.GetMessagesByUserId(c)
		}))
	}
	return r
}
