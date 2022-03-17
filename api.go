package main

import "encoding/json"
import "log"
import "net/http"

import "github.com/google/uuid"
import "github.com/gorilla/mux"

type API struct {
	router *mux.Router
	store  *Store
}

func NewAPI(storePath string) *API {
	api := &API{
		router: mux.NewRouter(),
		store:  NewStore(storePath),
	}

	api.router.HandleFunc("/template", api.createTemplate).Methods("POST").Headers("Content-Type", "application/json")
	api.router.HandleFunc("/template/{id}", api.streamTemplate).Methods("GET").Headers("Accept", "text/event-stream")
	api.router.HandleFunc("/template/{id}", api.getTemplate).Methods("GET")
	api.router.HandleFunc("/template/{id}", api.updateTemplate).Methods("PATCH").Headers("Content-Type", "application/json")

	return api
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api *API) createTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("createTemplate")

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	template := NewTemplate()
	err := dec.Decode(template)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	template.Id = uuid.NewString()

	if !template.IsValid() {
		http.Error(w, "invalid template", http.StatusBadRequest)
		return
	}

	err = api.store.Write(template)
	if err != nil {
		http.Error(w, "failed to write template", http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(template)
	if err != nil {
		http.Error(w, "failed to encode json", http.StatusInternalServerError)
		return
	}
}

func (api *API) streamTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("streamTemplate %s", mux.Vars(r))

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	closeChan := w.(http.CloseNotifier).CloseNotify()

	flusher.Flush()

	<-closeChan

	log.Printf("streamTemplate %s end", mux.Vars(r))
}

func (api *API) getTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("getTemplate %s", mux.Vars(r))
}

func (api *API) updateTemplate(w http.ResponseWriter, r *http.Request) {
	log.Printf("updateTemplate %s", mux.Vars(r))
}
