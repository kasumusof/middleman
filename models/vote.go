package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// Vote is used by pop to map your votes database table to your go code.
type Vote struct {
	ID        uuid.UUID `json:"id" db:"id"`
	VoterID   uuid.UUID `json:"voter_id" db:"voter_id"`
	VotedID   uuid.UUID `json:"voted_id" db:"voted_id"`
	BetID     uuid.UUID `json:"bet_id" db:"bet_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func checkExistenceOfVote(voted, voter *User, bet *Bet) bool {
	ubet := Vote{}
	nbet := Vote{}
	query := Tx.Where("voted_id = ? AND voter_id = ? AND bet_id = ?", voted.ID, voter.ID, bet.ID)

	if query.First(&ubet); nbet == ubet {
		return false
	}
	return true
}

func VoteBet(id, voted, voter string) (*Bet, error) {
	var uvoter, uvoted User
	var bet Bet
	query := Tx.Where("username = ?", voted)
	query.First(&uvoted)
	query = Tx.Where("username = ?", voter)
	query.First(&uvoter)
	query = Tx.Where("id = ?", id)
	query.First(&bet)

	if !checkExistenceOfUserInBet(&uvoter, &bet) || !checkExistenceOfUserInBet(&uvoted, &bet) {
		return nil, errors.New("Error shii happeneed here")

	}
	if checkExistenceOfVote(&uvoted, &uvoter, &bet) {
		return nil, errors.New("Vote already casted")
	}
	vote := Vote{BetID: bet.ID, VoterID: uvoter.ID, VotedID: uvoted.ID}
	if _, err := Tx.ValidateAndSave(&vote); err != nil {
		return nil, err
	}
	v, ok := bet.Votes[uvoted.Username].(float64)
	if !ok {
		log.Println(v, ok)
		return nil, errors.New("Unexpected error to float64")
	}

	bet.Votes[uvoted.Username] = v + 1
	Tx.ValidateAndSave(&bet)
	Tx.Load(&bet)
	return &bet, nil
}

func UnVoteBet(id, voted, voter string) (*Bet, error) {
	var uvoter, uvoted User
	var bet Bet
	query := Tx.Where("username = ?", voted)
	query.First(&uvoted)
	query = Tx.Where("username = ?", voter)
	query.First(&uvoter)
	query = Tx.Where("id = ?", id)
	query.First(&bet)

	if !checkExistenceOfVote(&uvoted, &uvoter, &bet) {
		return nil, errors.New("Vote not made yet")
	}

	var vote Vote
	query = Tx.Where("voted_id = ? AND voter_id = ? AND bet_id = ?", uvoted.ID, uvoter.ID, bet.ID)
	query.First(&vote)
	err = Tx.Destroy(&vote)
	if err != nil {
		return nil, err
	}
	v, ok := bet.Votes[uvoted.Username].(float64)
	if !ok {
		log.Println(v, ok)
		return nil, errors.New("Unexpected error to float64")
	}
	bet.Votes[uvoted.Username] = v - 1
	fmt.Printf("value: %v, type: %T", bet.Voters[uvoted.Username], bet.Voters[uvoted.Username])
	Tx.ValidateAndSave(&bet)
	Tx.Load(&bet)
	return &bet, nil
}

// String is not required by pop and may be deleted
func (v Vote) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Votes is not required by pop and may be deleted
type Votes []Vote

// String is not required by pop and may be deleted
func (v Votes) String() string {
	jv, _ := json.Marshal(v)
	return string(jv)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (v *Vote) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (v *Vote) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (v *Vote) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
