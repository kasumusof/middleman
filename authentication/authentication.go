package authentication

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	privateKeyPath = "keys/app.rsa"
	// pubKeyPath     = "keys/app.rsa.pub"
)

var VerifyKey, SignKey []byte

func init() {
	var err error

	SignKey, err = ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Error reading private key.")
		return
	}

	// VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	// if err != nil {
	// 	log.Fatal("Error reading public key.")
	// 	return
	// }
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token,omitempty"`
}

func GenerateToken(min int, username string) (*Token, error) {
	expirationTime := time.Now().Add(time.Duration(min) * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SignKey)
	if err != nil {
		return nil, err
	}

	return &Token{Token: tokenString}, nil
}

func VerifyToken(token string) bool {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return SignKey, nil
	})
	if err != nil {
		return false
	}

	if !tkn.Valid {
		return false
	}
	return true
}
