package main

import "encoding/json"
import "fmt"
import "log"
import "net/http"
import "time"

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

	api.router.HandleFunc("/template", returnError(jsonOutput(api.createTemplate))).Methods("POST").Headers("Content-Type", "application/json")
	api.router.HandleFunc("/template/{id}", returnError(api.streamTemplate)).Methods("GET").Headers("Accept", "text/event-stream")
	api.router.HandleFunc("/template/{id}", returnError(jsonOutput(api.getTemplate))).Methods("GET")
	api.router.HandleFunc("/template/{id}", returnError(jsonOutput(api.updateTemplate))).Methods("PATCH").Headers("Content-Type", "application/json")

	return api
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api *API) createTemplate(r *http.Request) (Object, string, int) {
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

func (api *API) streamTemplate(w http.ResponseWriter, r *http.Request) (string, int) {
	log.Printf("streamTemplate %s", mux.Vars(r))

	_, ok := w.(http.Flusher)
	if !ok {
		return "Streaming unsupported", http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	template := NewTemplate()
	template.Id = mux.Vars(r)["id"]

	err := api.store.Read(template)
	if err != nil {
		return fmt.Sprintf("Template %s not found", template.Id), http.StatusNotFound
	}

	writeEvent(w, template)

	closeChan := w.(http.CloseNotifier).CloseNotify()
	msgChan := api.bus.Subscribe(template)
	ticker := time.NewTicker(5 * time.Second)

	connected := true
	for connected {
		select {
		case <-closeChan:
			connected = false
		case msg := <-msgChan:
			writeEvent(w, msg)
		case <-ticker.C:
			writeEvent(w, NewHeartbeat())
		}
	}

	log.Printf("streamTemplate %s end", mux.Vars(r))

	return "", 0
}

func (api *API) getTemplate(r *http.Request) (Object, string, int) {
	log.Printf("getTemplate %s", mux.Vars(r))

	template := NewTemplate()
	template.Id = mux.Vars(r)["id"]

	err := api.store.Read(template)
	if err != nil {
		return nil, fmt.Sprintf("Template %s not found", template.Id), http.StatusNotFound
	}

	return template, "", 0
}

func (api *API) updateTemplate(r *http.Request) (Object, string, int) {
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

func readJson(r *http.Request, out Object) (string, int) {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(out)
	if err != nil {
		return fmt.Sprintf("Invalid JSON: %s", err), http.StatusBadRequest
	}

	return "", 0
}

func returnError(wrapped func(http.ResponseWriter, *http.Request) (string, int)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, code := wrapped(w, r)
		if code != 0 {
			http.Error(w, msg, code)
		}
	}
}

func jsonOutput(wrapped func(*http.Request) (Object, string, int)) func(http.ResponseWriter, *http.Request) (string, int) {
	return func(w http.ResponseWriter, r *http.Request) (string, int) {
		out, msg, code := wrapped(r)
		if code != 0 {
			return msg, code
		}

		enc := json.NewEncoder(w)
		err := enc.Encode(out)
		if err != nil {
			return fmt.Sprintf("Failed to encode JSON: %s", err), http.StatusInternalServerError
		}

		return "", 0
	}
}

func writeEvent(w http.ResponseWriter, in Object) (string, int) {
	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Sprintf("Failed to encode JSON: %s", err), http.StatusInternalServerError
	}

	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", in.GetType(), data)
	w.(http.Flusher).Flush()

	return "", 0
}
