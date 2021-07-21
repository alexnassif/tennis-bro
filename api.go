package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/alexnassif/tennis-bro/Auth"
	"github.com/alexnassif/tennis-bro/Models"
)

type LoginUser struct {
    Username string `json:"username"`
    Password string `json:"password"`
}



func HandleLogin(w http.ResponseWriter, r *http.Request) {

    var loginUser LoginUser

    // Try to decode the JSON request to a LoginUser
    err := json.NewDecoder(r.Body).Decode(&loginUser)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	var dbUser Models.User
    fmt.Println(loginUser.Username)
    // Find the user in the database by username
    userErr := Models.FindUserByUsername(loginUser.Username, &dbUser)
    if userErr != nil {
        returnErrorResponse(w)
        return
    }

    // Check if the passwords match
    ok, err := Auth.ComparePassword(loginUser.Password, dbUser.Password)

    if !ok || err != nil {
        returnErrorResponse(w)
        return
    }

    // Create a JWT
    token, err := Auth.CreateJWTToken(dbUser)

    if err != nil {
        returnErrorResponse(w)
        return
    }

    w.Write([]byte(token))

}

func returnErrorResponse(w http.ResponseWriter) {

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte("{\"status\": \"error\"}"))
}