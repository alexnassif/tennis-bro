package Auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnonUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (user *AnonUser) GetId() string {
	return user.Id
}

func (user *AnonUser) GetName() string {
	return user.Name
}

func AuthMiddleware(f gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {

		r := c.Request
		w := c.Writer
		token, tok := r.URL.Query()["bearer"]

		if tok && len(token) == 1 {
			user, err := ValidateToken(token[0])
            fmt.Print(user.GetName(), user.GetId())
			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)

			} else {
				c.Set("user", user)
				f(c)
			}

		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Please login or provide name"))
		}
	})
}
