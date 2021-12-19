package Controllers

import(
	"github.com/alexnassif/tennis-bro/Models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetMessagesByUserId(c *gin.Context) {
	//sender := c.Params.ByName("sender")
	recipient := c.Params.ByName("recipient")

	usertok, ok := c.Keys["user"].(Models.LoggedInUser)

	if !ok  {
		http.Error(c.Writer, "Forbidden", http.StatusForbidden)
		return
	}

	var messages []Models.Message
	err := Models.GetMessagesByUser(usertok.GetId(), recipient, &messages)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	   } else {
		c.JSON(http.StatusOK, messages)
	   }
}