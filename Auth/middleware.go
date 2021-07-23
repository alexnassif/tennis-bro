package Auth

import (
    "net/http"
    "github.com/google/uuid"
	"github.com/gin-gonic/gin"
)


type AnonUser struct {
    Id string `json:"id"`
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
        name, nok := r.URL.Query()["id"]

        if tok && len(token) == 1 {            
            user, err := ValidateToken(token[0])
            if err != nil {
                http.Error(w, "Forbidden", http.StatusForbidden)

            } else {
                c.Set("user", &user)
                f(c)
            }

        } else if nok && len(name) == 1 {
            // Continue with new Anon. user
            user := AnonUser{Id: uuid.New().String(), Name: name[0]}
            c.Set("user", &user)
            f(c)

        } else {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("Please login or provide name"))
        }
    })
}