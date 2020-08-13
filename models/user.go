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

// User is used by pop to map your users database table to your go code.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"-" db:"password"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	CreatedBets  Bets      `json:"created_bets,omitempty" has_many:"bets" order_by:"created_at asc"`
	BetsInvolved Bets      `json:"bets_involved,omitempty" many_to_many:"users_bets" order_by:"created_at asc"`
	Voters       Users     `json:"-" many_to_many:"votes" fk_id:"voted_id" order_by:"created_at asc"`
	Voteds       Users     `json:"-" many_to_many:"votes" fk_id:"voter_id" order_by:"created_at asc"`
}

type UserFormData struct {
	User
	Key string `json:"password"`
}

func CreateUser(user *User) (*User, error) {
	_, err = Tx.ValidateAndCreate(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUsers() (Users, error) {
	users := Users{}
	err = Tx.Eager().All(&users)
	if err != nil {
		return Users{}, err
	}
	return users, nil
}

func GetUser(username string) (*User, error) {

	user := User{}
	query := Tx.Where("username = ?", username)
	err = query.First(&user)
	if err != nil {
		return nil, err
	}
	Tx.Load(&user)
	return &user, nil
}

func GetUserByEmail(email string) (*User, error) {

	user := User{}
	query := Tx.Where("email = ?", email)
	err = query.First(&user)
	if err != nil {
		return nil, err
	}
	Tx.Load(&user)
	return &user, nil
}

func UpdateUserPassword(username, password string) (*User, error) {
	// fmt.Println("Update user")

	user := User{}
	query := Tx.Where("username = ?", username)
	err = query.First(&user)
	if err != nil {
		return nil, err
	}
	user.Password = password
	// Tx.ValidateAndCreate()
	// Tx.ValidateAndSave()
	// Tx.ValidateAndUpdate()
	_, err = Tx.ValidateAndSave(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUser(username string) (*User, error) {

	user := User{}
	query := Tx.Where("username = ?", username)
	err = query.First(&user)
	if err != nil {
		return nil, err
	}
	err = Tx.Destroy(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateBet(username string, bet *Bet) (*Bet, error) {
	var user User
	query := Tx.Where("username = ?", username)
	err := query.First(&user)
	if err != nil {
		return nil, err
	}
	bet.Votes = slices.Map{}
	bet.Votes[user.Username] = 0
	bet.Ongoing = true
	// bet.Voters[user.Username] = slices.String{}
	bet.Creator = &user
	_, err = Tx.ValidateAndCreate(bet)
	if err != nil {
		return nil, err
	}
	ubet := UsersBet{UserID: user.ID, BetID: bet.ID}
	_, err = Tx.ValidateAndSave(&ubet)
	if err != nil {
		return nil, err
	}
	Tx.Load(bet)
	return bet, nil
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Username, Name: "Username"},
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
