package Auth

import (
	"github.com/gin-gonic/gin"
 log "github.com/sirupsen/logrus"
)


func BasicAuth(c *gin.Context) {
	// Get the Basic Authentication credentials
	user, password, hasAuth := c.Request.BasicAuth()
	if hasAuth && user == "alexnassif" && password == "love" {
		log.WithFields(log.Fields{
			"user": user,
		}).Info("User authenticated")
	} else {
		c.Abort()
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		return
	}
}