package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	accesstokenBusinessLogic "github.com/konaro/line-notify-service/businesslogic/accesstoken"
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
	r := mux.NewRouter()

	r.HandleFunc("/oauth/line", lineAuthHandler).Methods("GET")

	api := r.PathPrefix("/api").Subrouter()

	api.Handle("/security", middleware.AuthHandler(http.HandlerFunc(securityHandler))).Methods("PATCH")
	api.Handle("/verify", middleware.AuthHandler(http.HandlerFunc(verifyHandler))).Methods("GET")
	api.Handle("/tokens", middleware.AuthHandler(http.HandlerFunc(tokenHandler))).Methods("GET")
	api.Handle("/tokens/{id}", middleware.AuthHandler(http.HandlerFunc(tokenHandler))).Methods("GET", "DELETE")
	api.Handle("/line-notify", middleware.AuthHandler(http.HandlerFunc(lineNotifyHandler))).Methods("POST")
	api.HandleFunc("/login", loginHandler).Methods("POST")
	api.HandleFunc("/line-callback", lineCallbackHandler).Methods("POST")

	http.ListenAndServe(":10000", r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
}

func securityHandler(w http.ResponseWriter, r *http.Request) {
	var entity model.ResetPassword

	err := json.NewDecoder(r.Body).Decode(&entity)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get claims account from context
	account := r.Context().Value("account").(string)

	// check old password correct
	if authBusinessLogic.CheckSecurity(account, entity.Password) {
		err = authBusinessLogic.UpdatePassword(entity.NewPassword)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// get tokens
	case "GET":
		res := accesstokenBusinessLogic.GetList(10, 0)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(model.Response{Data: res, Success: true})

	// delete token
	case "DELETE":
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = accesstokenBusinessLogic.Delete(id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func lineAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := linenotify.GetAuthorizeUrl(clientId, host+"/line-callback")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func lineCallbackHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.FormValue("code")

	token := linenotify.GetAccessToken(code, clientId, secret, host+"/line-callback")

	response := &model.TokenResponse{}
	json.Unmarshal(token, response)

	// add token to storage
	accesstokenBusinessLogic.Add(response.AccessToken)
}

func lineNotifyHandler(w http.ResponseWriter, r *http.Request) {
	var notify model.Notify

	err := json.NewDecoder(r.Body).Decode(&notify)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get all tokens
	tokens := accesstokenBusinessLogic.GetAllTokens()

	requests := make(chan []byte)

	var wg sync.WaitGroup
	wg.Add(len(tokens))

	for _, token := range tokens {
		go func(notify model.Notify, token string) {
			defer wg.Done()
			res, _ := linenotify.PushNotification(notify, token)
			requests <- res
		}(notify, token)
	}

	go func() {
		for response := range requests {
			fmt.Println(string(response))
		}
	}()

	wg.Wait()
}
