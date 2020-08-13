package routes

import (
	"net/http"

	"github.com/kasumusof/sodhero/api/controllers"
)

var betRoutes = []Route{
	{
		URL:     "/bets",
		Method:  http.MethodGet,
		Handler: controllers.GetBets,
	},
	{
		URL:     "/bets/{id}",
		Method:  http.MethodGet,
		Handler: controllers.GetBet,
	},
	{
		URL:     "/bets",
		Method:  http.MethodPost,
		Handler: controllers.CreateBet,
	},
	{
		URL:     "/bets/{id}",
		Method:  http.MethodPut,
		Handler: controllers.UpdateBet,
	},
	{
		URL:     "/bets/{id}",
		Method:  http.MethodDelete,
		Handler: controllers.DeleteBet,
	},
	{
		URL:     "/bets/{id}/users",
		Method:  http.MethodGet,
		Handler: controllers.GetBetUsers,
	},
	{
		URL:     "/bets/{id}/add_user/{username}",
		Method:  http.MethodPut,
		Handler: controllers.AddUserToBet,
	},
	{
		URL:     "/bets/{id}/delete_user/{username}",
		Method:  http.MethodDelete,
		Handler: controllers.DeleteUserFromBet,
	},
	{
		URL:     "/bets/{id}/{user1}/vote/{user2}",
		Method:  http.MethodPut,
		Handler: controllers.VoteBet,
	},
	{
		URL:     "/bets/{id}/{user1}/unvote/{user2}",
		Method:  http.MethodPut,
		Handler: controllers.UnVoteBet,
	},
}
