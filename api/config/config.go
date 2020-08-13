package config

import "os"

var (
	// PORT to listen on
	PORT = ""
)

func init() {
	PORT = os.Getenv("PORT")
}
