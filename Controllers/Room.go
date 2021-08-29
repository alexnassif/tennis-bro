package Controllers

import(
	"github.com/alexnassif/tennis-bro/Models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetRoomsForUser(c *gin.Context) {
	id := c.Params.ByName("id")

	usertok, ok := c.Keys["user"].(Models.LoggedInUser)
	if !ok {
		return
	}

	if usertok.GetId() != id {
		return
	}
	var privateRooms []Models.Room
	err := Models.GetRoomsByUsers(id, &privateRooms)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	   } else {
		c.JSON(http.StatusOK, privateRooms)
	   }
}