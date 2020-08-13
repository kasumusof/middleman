package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/kasumusof/sodhero/models"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users models.Users
	users, err := models.GetUsers()
	if err != nil {
		log.Println(err)
		makeCache(r, makeResponse(&w, "Internal Server Error", false, http.StatusInternalServerError))
		return
	}
	makeCache(r, makeResponse(&w, users, true, http.StatusOK))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.GetUser(username)
	if err != nil {
		makeCache(r, makeResponse(&w, "User does not exist", false, http.StatusBadRequest))
		return
	}
	makeCache(r, makeResponse(&w, user, true, http.StatusOK))
}

// func GetUserBets(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	username := vars["username"]
// 	user, err := models.GetUser(username)
// 	if err != nil {
// 		makeCache(r, makeResponse(&w, "User does not exist", false, http.StatusBadRequest))
// 		return
// 	}
// 	makeCache(r, makeResponse(&w, user.BetsInvolved, true, http.StatusOK))
// }

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	makeResponse(&w, "Update User: Not implementted", false, http.StatusOK)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := models.DeleteUser(username)
	if err != nil {
		makeResponse(&w, "User does not exist", false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, user, true, http.StatusOK)

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	uCtx, _ := context.Get(r, "foo").(*models.User)
	u, err := models.CreateUser(uCtx)
	if err != nil || u.ID == uuid.Nil {
		makeResponse(&w, "Error in creating user", false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, u, true, http.StatusOK)
}

func CreateBet(w http.ResponseWriter, r *http.Request) {
	var bet models.Bet
	// var response models.Response
	username := mux.Vars(r)["username"]

	err := json.NewDecoder(r.Body).Decode(&bet)
	if err != nil {
		log.Println("Error:", err)
	}
	_, err = models.CreateBet(username, &bet)
	if err != nil {
		log.Println("Error:", err)
		// w.WriteHeader(http.StatusBadRequest)
		// response.OK = false
		// response.Result = "Error in creating bet"
		// json.NewEncoder(w).Encode(response)
		makeResponse(&w, "error in creating bet", false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, bet, true, http.StatusOK)
}
