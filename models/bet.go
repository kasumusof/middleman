package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/pop/slices"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// Bet is used by pop to map your bets database table to your go code.
type Bet struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	BetName      string     `json:"bet_name" db:"bet_name"`
	Ongoing      bool       `json:"ongoing" db:"ongoing"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	CreatorID    uuid.UUID  `json:"-" db:"user_id"`
	Creator      *User      `json:"creator,omitempty" belongs_to:"user"`
	Participants Users      `json:"users,omitempty" many_to_many:"users_bets"`
	VoteBet      Bets       `json:"-" many_to_many:"votes" fk_id:"bet_id"`
	Votes        slices.Map `json:"votes" db:"votes"`
	Voters       slices.Map `json:"-" db:"voters"`
}

func GetBets() (Bets, error) {

	bets := Bets{}
	err = Tx.Eager().All(&bets)
	if err != nil {
		return Bets{}, err
	}
	return bets, nil
}

func GetBet(id string) (*Bet, error) {

	bet := Bet{}
	err = Tx.Find(&bet, id)
	if err != nil {
		return nil, err
	}
	Tx.Load(&bet)
	return &bet, nil
}

func UpdateBet(id, name string) (*Bet, error) {

	bet := Bet{}
	err = Tx.Find(&bet, id)
	if err != nil {
		return nil, err
	}
	bet.BetName = name
	// Tx.ValidateAndCreate()
	// Tx.ValidateAndSave()
	// Tx.ValidateAndUpdate()
	_, err = Tx.ValidateAndSave(&bet)
	if err != nil {
		return nil, err
	}
	Tx.Load(&bet)
	return &bet, nil

}

func DeleteBet(id string) (*Bet, error) {

	bet := Bet{}
	err = Tx.Find(&bet, id)
	if err != nil {
		return nil, err
	}
	Tx.Load(&bet)
	err = Tx.Destroy(&bet)
	if err != nil {
		return nil, err
	}
	return &bet, nil
}

// String is not required by pop and may be deleted
func (b Bet) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Bets is not required by pop and may be deleted
type Bets []Bet

// String is not required by pop and may be deleted
func (b Bets) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *Bet) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: b.BetName, Name: "BetName"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *Bet) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *Bet) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
