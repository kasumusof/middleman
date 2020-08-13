package router

import (
	"net/http"
	"os"
	"path/filepath"

	log "github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/kasumusof/sodhero/api/router/routes"
	auth "github.com/kasumusof/sodhero/authentication"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

var a = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
var logger = log.With(a, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

// New creates new router object for api pakcage
func New() *mux.Router {
	r := setupRouter(newRouter())
	return r
}

func newRouter() *mux.Router {
	return mux.NewRouter()
}

func setupRouter(r *mux.Router) *mux.Router {
	// route := routes.MainRoutes()
	// r.Handle("/", http.StripPrefix("/",http.FileServer(http.Dir("./gasbet/build"))))
	// r.HandleFunc("/", auth.EnableCors(controllers.Home("/home/yugoslavman/Documents/go/src/github.com/kasumusof/NewToTest/gasbet/build/index.html")))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer((http.Dir("gasbet/build/static/")))))

	api := r.PathPrefix("/api/v1/").Subrouter()
	for _, route := range routes.ApiRoutes() {
		// // log.Println(route)
		// if route.URL == "/" {
		// 	continue
		// }
		if route.URL == "/getToken" {
			api.HandleFunc(route.URL, auth.Jsonify(auth.EnableCors(auth.Cacher(route.Handler)))).Methods(route.Method)
			continue
		}

		if route.URL == "/users" && route.Method == http.MethodPost {
			api.HandleFunc(route.URL, auth.Jsonify(auth.EnableCors(auth.ValidateRegister(route.Handler)))).Methods(route.Method)
			continue
		}

		// if route.Method == http.MethodPost {
		// 	api.HandleFunc(route.URL, auth.Jsonify(auth.EnableCors(route.Handler))).Methods(route.Method)
		// 	continue
		// }

		api.HandleFunc(route.URL, auth.Jsonify(auth.EnableCors(auth.Authenticate(auth.Cacher(route.Handler))))).Methods(route.Method)
	}

	// cor := cors.AllowAll()
	r.Use(auth.LoggingMiddlewareReal(logger))
	// fs := http.FileServer(http.Dir("static"))
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static", fs))
	// r.Use(auth.Authenticate)
	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)
	return r
}
