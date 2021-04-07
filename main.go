package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	authBusinessLogic "github.com/konaro/line-notify-service/businesslogic/auth"
	"github.com/konaro/line-notify-service/client/linenotify"
	"github.com/konaro/line-notify-service/jwtutil"
	"github.com/konaro/line-notify-service/middleware"
	"github.com/konaro/line-notify-service/model"
)

var clientId = os.Getenv("clientId")
var secret = os.Getenv("secret")
var host = os.Getenv("host")

func main() {
	handleRequests()
}

func handleRequests() {
	mux := http.NewServeMux()

	mux.Handle("/security", middleware.AuthHandler(http.HandlerFunc(security)))
	mux.Handle("/verify", middleware.AuthHandler(http.HandlerFunc(verify)))
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/line-auth", lineAuth)
	mux.HandleFunc("/callback", callback)
	log.Fatal(http.ListenAndServe(":10000", mux))
}

// login handler
func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var login model.Login

		err := json.NewDecoder(r.Body).Decode(&login)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := authBusinessLogic.CheckSecurity(login.Account, login.Password)

		if result {
			// generate token
			token, err := jwtutil.GenerateToken(login.Account)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(model.Response{Data: token, Success: true})
			return
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func security(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "PATCH":
		var entity model.ResetPassword

		err := json.NewDecoder(r.Body).Decode(&entity)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get claims account from context
		account := r.Context().Value("account").(string)

		// check old password valid
		if authBusinessLogic.CheckSecurity(account, entity.Password) {
			err = authBusinessLogic.UpdatePassword(entity.NewPassword)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			w.WriteHeader(http.StatusForbidden)
		}

		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func verify(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func lineAuth(w http.ResponseWriter, r *http.Request) {
	url := linenotify.GetAuthorizeUrl(clientId, host+"/callback")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func callback(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.FormValue("code")

	token := linenotify.GetAccessToken(code, clientId, secret, host+"/callback")

	response := &model.TokenResponse{}
	json.Unmarshal(token, response)

	fmt.Printf("accesstoken: %v", response.AccessToken)
}
