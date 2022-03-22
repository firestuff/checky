package main

// import "fmt"
// import "github.com/google/uuid"

import "flag"
import "log"
import "net/http"

var bindFlag = flag.String("listen", "[::]:8100", "host:port to listen on")
var storeFlag = flag.String("store", "data", "data store path")

func main() {
	flag.Parse()

	mux := http.NewServeMux()

	api, err := NewAPI(*storeFlag)
	if err != nil {
		log.Fatal(err)
	}

	mux.Handle("/api/", http.StripPrefix("/api", api))

	srv := &http.Server{
		Addr:    *bindFlag,
		Handler: mux,
	}

	log.Printf("listening on %s", *bindFlag)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
