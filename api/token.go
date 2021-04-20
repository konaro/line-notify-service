package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/konaro/line-notify-service/businesslogic/accesstoken"
	"github.com/konaro/line-notify-service/client/linenotify"
	"github.com/konaro/line-notify-service/model"
)

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	res := accesstoken.GetList(10, 0)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.Response{Data: res, Success: true})
}

func TokenRevokeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	token := accesstoken.Get(id)

	resp, err := linenotify.RevokeAccessToken(token.Token)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	response := &model.RevokeTokenResponse{}
	json.Unmarshal(resp, response)

	switch response.Status {
	case 200:
		err = accesstoken.Delete(id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusUnauthorized)
	}
}
