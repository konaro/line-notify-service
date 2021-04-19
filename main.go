package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/konaro/line-notify-service/api"
	"github.com/konaro/line-notify-service/middleware"
)

func main() {
	handleRequests()
}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/oauth/line", api.LineAuthHandler).Methods("GET")

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.Handle("/security", middleware.AuthHandler(http.HandlerFunc(api.SecurityHandler))).Methods("PATCH")
	apiRouter.Handle("/verify", middleware.AuthHandler(http.HandlerFunc(api.VerifyHandler))).Methods("GET")
	apiRouter.Handle("/tokens", middleware.AuthHandler(http.HandlerFunc(api.TokenHandler))).Methods("GET")
	apiRouter.Handle("/tokens/{id}", middleware.AuthHandler(http.HandlerFunc(api.TokenHandler))).Methods("GET")
	apiRouter.Handle("/tokens/revoke/{id}", middleware.AuthHandler(http.HandlerFunc(api.TokenRevokeHandler))).Methods("DELETE")
	apiRouter.Handle("/line-notify", middleware.AuthHandler(http.HandlerFunc(api.LineNotifyHandler))).Methods("POST")
	apiRouter.HandleFunc("/login", api.LoginHandler).Methods("POST")
	apiRouter.HandleFunc("/line-callback", api.LineCallbackHandler).Methods("POST")

	http.ListenAndServe(":10000", r)
}
