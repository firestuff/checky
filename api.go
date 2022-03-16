package main

import "log"
import "net/http"

import "github.com/gorilla/mux"

type API struct {
	router *mux.Router
}

func NewAPI() *API {
	api := &API{
		router: mux.NewRouter(),
	}

	api.router.HandleFunc("/template", api.createTemplate).Methods("POST")
	api.router.HandleFunc("/template/{id}", api.getTemplate).Methods("GET")
	api.router.HandleFunc("/template/{id}", api.updateTemplate).Methods("PATCH")

	return api
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api *API) createTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("createTemplate")
}

func (api *API) getTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("getTemplate %s", mux.Vars(r))
}

func (api *API) updateTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("updateTemplate %s", mux.Vars(r))
}
