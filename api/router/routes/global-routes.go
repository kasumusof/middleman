package routes

import (
	"net/http"

	"github.com/kasumusof/sodhero/api/controllers"
)

var mainRoute = Route{
	URL:     "/",
	Handler: controllers.Home("index.html"),
	Method:  http.MethodGet,
}

var unClassifieedRoutes = Routes{
	{
		URL:     "/getToken",
		Handler: controllers.Signin,
		Method:  http.MethodPost,
	},
}
