package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/konaro/line-notify-service/businesslogic/accesstoken"
	"github.com/konaro/line-notify-service/client/linenotify"
	"github.com/konaro/line-notify-service/model"
)

var clientId = os.Getenv("clientId")
var secret = os.Getenv("secret")
var host = os.Getenv("host")

func LineAuthHandler(w http.ResponseWriter, r *http.Request) {
	url := linenotify.GetAuthorizeUrl(clientId, host+"/api/line-callback")
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func LineCallbackHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.FormValue("code")

	token := linenotify.GetAccessToken(code, clientId, secret, host+"/api/line-callback")

	response := &model.TokenResponse{}
	json.Unmarshal(token, response)

	// add token to storage
	accesstoken.Add(response.AccessToken)
}

func LineNotifyHandler(w http.ResponseWriter, r *http.Request) {
	var notify model.Notify

	err := json.NewDecoder(r.Body).Decode(&notify)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get all tokens
	tokens := accesstoken.GetAllTokens()

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
