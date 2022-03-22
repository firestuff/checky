package main

import "fmt"
import "net/http"

import "github.com/firestuff/storebus"

type API struct {
	api *storebus.API
}

func NewAPI(root string) (*API, error) {
	api := &API{}

	var err error
	api.api, err = storebus.NewAPI(
		root,
		&storebus.APIConfig{
			Factory:   factory,
			Update:    update,
			MayCreate: mayCreate,
			MayUpdate: mayUpdate,
			MayRead:   mayRead,
		},
	)

	if err != nil {
		return nil, err
	}

	return api, nil
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.api.ServeHTTP(w, r)
}

func factory(t string) (storebus.Object, error) {
	switch t {

	case "template":
		return NewTemplate(), nil

	default:
		return nil, fmt.Errorf("Unsupported type: %s", t)

	}
}

func update(obj storebus.Object, patch storebus.Object) error {
	switch o := obj.(type) {

	case *Template:
		p := patch.(*Template)

		if p.Title != "" {
			o.Title = p.Title
		}

		return nil

	default:
		return fmt.Errorf("Unsupported type: %s", obj.GetType())

	}
}

func mayCreate(obj storebus.Object, r *http.Request) error {
	return nil
}

func mayUpdate(obj storebus.Object, patch storebus.Object, r *http.Request) error {
	return nil
}

func mayRead(obj storebus.Object, r *http.Request) error {
	return nil
}
