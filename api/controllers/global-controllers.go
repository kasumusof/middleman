package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/kasumusof/sodhero/authentication"
	"github.com/kasumusof/sodhero/models"
)

type User struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

// Home endpoint
func Home(entryPoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entryPoint)
	}

	return http.HandlerFunc(fn)
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		makeResponse(&w, "error decoding", false, http.StatusBadRequest)
		return
	}
	if user.Password == "" || user.Username == "" {
		makeResponse(&w, "Empty Username", false, http.StatusBadRequest)

		return
	}
	u, err := models.GetUser(user.Username)
	if err != nil {
		makeResponse(&w, "Invalid Username", false, http.StatusBadRequest)
		return
	}
	if u.Password != user.Password {
		makeResponse(&w, "Invalid Password", false, http.StatusForbidden)
		return
	}
	token, err := authentication.GenerateToken(25, user.Username)
	if err != nil {
		makeResponse(&w, "Internal server Error", false, http.StatusBadRequest)
		return
	}
	toSend := struct {
		Username string `json:"username,omitempty"`
		Token    string `json:"token,omitempty"`
	}{
		Username: u.Username,
		Token: token.Token,
	}
	makeResponse(&w, toSend, true, http.StatusOK)
}

// func enableCors(w *http.ResponseWriter) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// }

func makeCache(r *http.Request, data interface{}) {
	body, _ := json.Marshal(data)
	context.Set(r, "tocache", body)
}

func makeResponse(w *http.ResponseWriter, result interface{}, ok bool, status int) models.Response {
	var response models.Response
	response.OK = ok
	response.Result = result
	(*w).WriteHeader(status)
	if err := json.NewEncoder(*w).Encode(response); err != nil {
		fmt.Println(err)
	}
	return response
}
