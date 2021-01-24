package Routes

import (
	"github.com/alexnassif/tennis-bro/Controllers"
	"github.com/gin-gonic/gin"
	"github.com/alexnassif/tennis-bro/Auth"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userGroup := r.Group("/user-api")
	{
		userGroup.GET("user", Auth.BasicAuth, Controllers.GetUsers)
		userGroup.POST("user", Controllers.CreateUser)
		userGroup.GET("user/:id", Controllers.GetUserByID)
		userGroup.PUT("user/:id", Controllers.UpdateUser)
		userGroup.DELETE("user/:id", Controllers.DeleteUser)
	}
	return r
}

