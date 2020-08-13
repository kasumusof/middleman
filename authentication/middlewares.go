package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/context"
	"github.com/kasumusof/sodhero/models"
)

func Authenticate(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer, ok := r.Header["Authorization"]
		if !ok {
			http.Error(w, "no authentication", http.StatusForbidden)
			return
		}
		// fmt.Println("1:", bearer)
		// fmt.Println("2:", bearer[0])
		// fmt.Println("3:", bearer[0][7:])
		// fmt.Println("4:", len(bearer))
		if ok := VerifyToken(bearer[0][7:]); !ok {
			http.Error(w, "unauthorized", http.StatusForbidden)
			return
		}

		next(w, r)
	}

}

func EnableCors(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("AllowCredentials","true")
		next(w, r)
	}
}

func Jsonify(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		next(w, r)
	}
}

func ValidateRegister(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var userForm models.UserFormData
		var response models.Response

		err := json.NewDecoder(r.Body).Decode(&userForm)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response.OK = false
			response.Result = "Error in json passed from middleware"
			json.NewEncoder(w).Encode(response)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		re := regexp.MustCompile(`^\w+$`)
		if !re.MatchString(userForm.Username) || len(userForm.Username) < 6 {
			response.OK = false
			response.Result = "invalid username"
			json.NewEncoder(w).Encode(response)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		re = regexp.MustCompile(`^(\w+\.?\w*)\+?(\w*)@(\w+)(\.\w+)+$`)
		if !re.MatchString(userForm.Email) {
			response.OK = false
			response.Result = "invalid email"
			json.NewEncoder(w).Encode(response)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		re = regexp.MustCompile(`^.*$`)
		if !re.MatchString(userForm.Key) || len(userForm.Key) < 8 {
			response.OK = false
			response.Result = "invalid Password"
			json.NewEncoder(w).Encode(response)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		user, _ := models.GetUser(userForm.Username)
		if user != nil {
			response.OK = false
			response.Result = "Username Taken"
			json.NewEncoder(w).Encode(response)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		user, _ = models.GetUserByEmail(userForm.Email)

		if user != nil {
			response.OK = false
			response.Result = "Email Taken"
			json.NewEncoder(w).Encode(response)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		user = &models.User{
			Username: userForm.Username,
			Password: userForm.Key,
			Email:    userForm.Email,
		}
		context.Set(r, "foo", user)
		next(w, r)
	}
}

func Cacher(next func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var response interface{}
		cache := NewCache("localhost:6379", "", 0)
		url := r.RequestURI
		if r.Method == http.MethodGet {
			if cachedItem, err := cache.Get(url); err == nil {
				fmt.Println("serving form cache")
				err = json.Unmarshal([]byte(cachedItem), &response)
				if err == nil {
					json.NewEncoder(w).Encode(response)
					return
				}
				fmt.Println("error reading cache", err, "now serving from db")
			}
			next(w, r)
			data := context.Get(r, "tocache")
			if err := cache.Set(r.RequestURI, data, 10*time.Second); err != nil {
				fmt.Println("error setting cache for request:", r.RequestURI)
			}
			return
		}
		next(w, r)
	}

}
