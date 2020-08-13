package api

import (
	"log"
	"net/http"

	"github.com/kasumusof/sodhero/api/config"
	"github.com/kasumusof/sodhero/api/router"
	"github.com/rs/cors"
)

// Run to run the api server
func Run() {
	serve()
}

func serve() {
	port := config.PORT
	if port == "" {
		log.Fatal("PORT needs to be set")
	}
	r := router.New()
	log.Println("listening on port", port)
	c := cors.New(cors.Options{
		// AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":"+port, handler))

}
