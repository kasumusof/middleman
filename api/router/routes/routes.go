package routes

import "net/http"

// Route the singular Routes
type Route struct {
	URL string
	// Handler interface{}
	Handler func(http.ResponseWriter, *http.Request)
	Method  string
}

// Routes the plural Route
type Routes []Route

// AvailableRoutes give available routes
func ApiRoutes() (routes Routes) {
	routes = append(unClassifieedRoutes, userRoutes...)
	routes = append(routes, betRoutes...)
	// routes = append(routes, fileRoutes...)
	return
}

// MainRoutes give available routes
func MainRoutes() (route Route) {
	route = mainRoute
	return
}
