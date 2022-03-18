package main

import "encoding/json"
import "fmt"
import "log"
import "net/http"

import "github.com/google/uuid"
import "github.com/gorilla/mux"

type API struct {
	router *mux.Router
	store  *Store
	bus    *Bus
}

func NewAPI(storePath string) *API {
	api := &API{
		router: mux.NewRouter(),
		store:  NewStore(storePath),
		bus:    NewBus(),
	}

	api.router.HandleFunc("/template", jsonOutput(api.createTemplate)).Methods("POST").Headers("Content-Type", "application/json")
	api.router.HandleFunc("/template/{id}", api.streamTemplate).Methods("GET").Headers("Accept", "text/event-stream")
	api.router.HandleFunc("/template/{id}", jsonOutput(api.getTemplate)).Methods("GET")
	api.router.HandleFunc("/template/{id}", jsonOutput(api.updateTemplate)).Methods("PATCH").Headers("Content-Type", "application/json")

	return api
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api *API) createTemplate(r *http.Request) (interface{}, string, int) {
	log.Printf("createTemplate")

	template := NewTemplate()
	msg, code := readJson(r, template)
	if code != 0 {
		return nil, msg, code
	}

	template.Id = uuid.NewString()

	if !template.IsValid() {
		return nil, "Invalid template", http.StatusBadRequest
	}

	err := api.store.Write(template)
	if err != nil {
		return nil, "Failed to write template", http.StatusInternalServerError
	}

	return template, "", 0
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

func (api *API) getTemplate(r *http.Request) (interface{}, string, int) {
	log.Printf("getTemplate %s", mux.Vars(r))

	template := NewTemplate()
	template.Id = mux.Vars(r)["id"]

	err := api.store.Read(template)
	if err != nil {
		return nil, fmt.Sprintf("Template %s not found", template.Id), http.StatusNotFound
	}

	return template, "", 0
}

func (api *API) updateTemplate(r *http.Request) (interface{}, string, int) {
	log.Printf("updateTemplate %s", mux.Vars(r))

	patch := NewTemplate()

	msg, code := readJson(r, patch)
	if code != 0 {
		return nil, msg, code
	}

	template := NewTemplate()
	template.Id = mux.Vars(r)["id"]

	err := api.store.Read(template)
	if err != nil {
		return nil, fmt.Sprintf("Template %s not found", template.Id), http.StatusNotFound
	}

	if patch.Title != "" {
		template.Title = patch.Title
	}

	if !template.IsValid() {
		return nil, "Invalid template", http.StatusBadRequest
	}

	err = api.store.Write(template)
	if err != nil {
		return nil, "Failed to write template", http.StatusInternalServerError
	}

	api.bus.Announce(template)

	return template, "", 0

}

func readJson(r *http.Request, out interface{}) (string, int) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(out)
	if err != nil {
		return fmt.Sprintf("Invalid JSON: %s", err), http.StatusBadRequest
	}

	return "", 0
}

func jsonOutput(wrapped func(*http.Request) (interface{}, string, int)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		out, msg, code := wrapped(r)
		if code != 0 {
			http.Error(w, msg, code)
			return
		}

		enc := json.NewEncoder(w)
		err := enc.Encode(out)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode JSON: %s", err), http.StatusInternalServerError)
		}
	}
}
