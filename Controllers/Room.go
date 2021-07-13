package Controllers

import(
	"github.com/alexnassif/tennis-bro/Models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetRoomsForUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var privateRooms []Models.Room
	err := Models.GetRoomsByUsers(id, &privateRooms)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	   } else {
		c.JSON(http.StatusOK, privateRooms)
	   }
}