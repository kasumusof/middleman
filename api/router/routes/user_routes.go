package routes

import (
	"net/http"

	"github.com/kasumusof/sodhero/api/controllers"
)

var userRoutes = []Route{
	{
		URL:     "/users",
		Method:  http.MethodGet,
		Handler: controllers.GetUsers,
	},
	{
		URL:     "/users/{username}",
		Method:  http.MethodGet,
		Handler: controllers.GetUser,
	},
	{
		URL:     "/users",
		Method:  http.MethodPost,
		Handler: controllers.CreateUser,
	},
	{
		URL:     "/users/{username}",
		Method:  http.MethodPut,
		Handler: controllers.UpdateUser,
	},
	{
		URL:     "/users/{username}",
		Method:  http.MethodDelete,
		Handler: controllers.DeleteUser,
	},
	{
		URL:     "/users/{username}/create_bet",
		Method:  http.MethodPost,
		Handler: controllers.CreateBet,
	},
}
