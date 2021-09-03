package Controllers

import (
	"net/http"

	"github.com/alexnassif/tennis-bro/Models"
	"github.com/gin-gonic/gin"
)

func GetRoomsForUser(c *gin.Context) {
	id := c.Params.ByName("id")

	usertok, ok := c.Keys["user"].(Models.LoggedInUser)
	if !ok || usertok.GetId() != id {
		http.Error(c.Writer, "Forbidden", http.StatusForbidden)
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
