package models

import (
	"log"

	"github.com/gobuffalo/pop"
)

var Tx *pop.Connection
var err error

func init() {
	var err error
	Tx, err = pop.Connect("development")
	if err != nil {
		log.Panic("error connecting to database", err)
	}

}

type Response struct {
	OK     bool        `json:"ok"`
	Result interface{} `json:"result"`
}
