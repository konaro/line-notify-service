package api

import (
	"encoding/json"
	"net/http"

	"github.com/konaro/line-notify-service/businesslogic/auth"
	"github.com/konaro/line-notify-service/jwtutil"
	"github.com/konaro/line-notify-service/model"
)

func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func SecurityHandler(w http.ResponseWriter, r *http.Request) {
	var entity model.ResetPassword

	err := json.NewDecoder(r.Body).Decode(&entity)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get claims account from context
	account := r.Context().Value("account").(string)

	// check old password correct
	if auth.CheckSecurity(account, entity.Password) {
		err = auth.UpdatePassword(entity.NewPassword)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var login model.Login

	err := json.NewDecoder(r.Body).Decode(&login)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := auth.CheckSecurity(login.Account, login.Password)

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
