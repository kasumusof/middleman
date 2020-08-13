package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kasumusof/sodhero/models"
)

func GetBets(w http.ResponseWriter, r *http.Request) {

	var bets models.Bets
	bets, err := models.GetBets()
	if err != nil {
		log.Println(err)
		makeCache(r, makeResponse(&w, "An error occured", false, http.StatusInternalServerError))
		return
	}
	makeCache(r, makeResponse(&w, bets, true, http.StatusOK))
}

func GetBet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bet, err := models.GetBet(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		makeCache(r, makeResponse(&w, "bet does not exist", false, http.StatusBadRequest))
		return
	}
	makeCache(r, makeResponse(&w, bet, true, http.StatusOK))

}

func GetBetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bet, err := models.GetBet(id)
	if err != nil {
		makeCache(r, makeResponse(&w, "bet does not exist", false, http.StatusBadRequest))
		return
	}
	makeCache(r, makeResponse(&w, bet.Participants, true, http.StatusOK))
}

func UpdateBet(w http.ResponseWriter, r *http.Request) {
	makeCache(r, makeResponse(&w, "Update Bet: Not implementted", false, http.StatusOK))
}

func DeleteBet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	bet, err := models.DeleteBet(id)
	if err != nil {
		makeResponse(&w, "bet does not exist", false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, bet, true, http.StatusOK)
}

func AddUserToBet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	id := vars["id"]
	bet, err := models.AddUserToBet(username, id)
	if err != nil {
		log.Println("Error Adding User to Bet:", err)
		makeResponse(&w, false, false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, bet, true, http.StatusOK)
}

func DeleteUserFromBet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	id := vars["id"]
	bet, err := models.DeleteUserFromBet(username, id)
	if err != nil {
		makeResponse(&w, false, false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, bet, true, http.StatusOK)
}

func VoteBet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	user1 := vars["user1"]
	user2 := vars["user2"]
	bet, err := models.VoteBet(id, user2, user1)
	if err != nil {
		makeResponse(&w, err, false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, bet, true, http.StatusOK)
}

func UnVoteBet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	user1 := vars["user1"]
	user2 := vars["user2"]
	bet, err := models.UnVoteBet(id, user2, user1)
	if err != nil {
		log.Println(err)
		makeResponse(&w, err, false, http.StatusBadRequest)
		return
	}
	makeResponse(&w, bet, true, http.StatusOK)

}
