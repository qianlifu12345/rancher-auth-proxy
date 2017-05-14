package service

import "github.com/gorilla/mux"

//NewRouter creates and configures a mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/v1-auth-filter/validateAuthToken", ValidationHandler).Methods("POST")
	return router
}
