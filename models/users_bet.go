package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// UsersBet is used by pop to map your users_bets database table to your go code.
type UsersBet struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	BetID     uuid.UUID `json:"bet_id" db:"bet_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func checkExistenceOfUserInBet(user *User, bet *Bet) bool {
	ubet := UsersBet{}
	bbet := UsersBet{}
	query := Tx.Where("user_id = ? AND bet_id = ?", user.ID, bet.ID)

	if query.First(&ubet); ubet == bbet {
		return false
	}
	return true
}

// AddUserToBet func
func AddUserToBet(username, id string) (*Bet, error) {
	user, err := GetUser(username)
	if err != nil {
		return nil, err
	}
	bet, err := GetBet(id)
	if err != nil {
		return nil, err
	}
	if checkExistenceOfUserInBet(user, bet) {
		return nil, err
	}
	ubet := UsersBet{UserID: user.ID, BetID: bet.ID}
	_, err = Tx.ValidateAndSave(&ubet)
	if err != nil {
		return nil, err
	}
	bet.Votes[user.Username] = 0
	bet.Voters[user.Username] = []string{}

	_, err = Tx.ValidateAndSave(bet)
	if err != nil {
		return nil, err
	}
	Tx.Load(bet)
	return bet, nil
}

func DeleteUserFromBet(username, id string) (*Bet, error) {
	user, err := GetUser(username)
	if err != nil {
		return nil, err
	}
	bet, err := GetBet(id)
	if err != nil {
		return nil, err
	}
	if !checkExistenceOfUserInBet(user, bet) {
		return nil, err
	}

	ubet := UsersBet{UserID: user.ID, BetID: bet.ID}
	query := Tx.Where("user_id = ? AND bet_id = ?", user.ID, bet.ID)
	query.First(&ubet)
	err = Tx.Destroy(&ubet)
	if err != nil {
		return nil, err
	}
	delete(bet.Votes, user.Username)
	delete(bet.Voters, user.Username)
	_, err = Tx.ValidateAndSave(bet)
	if err != nil {
		return nil, err
	}
	Tx.Load(bet)
	return bet, nil
}

// String is not required by pop and may be deleted
func (u UsersBet) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// UsersBets is not required by pop and may be deleted
type UsersBets []UsersBet

// String is not required by pop and may be deleted
func (u UsersBets) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *UsersBet) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *UsersBet) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *UsersBet) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
