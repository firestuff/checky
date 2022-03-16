package main

// import "fmt"
// import "github.com/google/uuid"

import "flag"
import "log"
import "net/http"

var bindFlag = flag.String("listen", "[::]:8100", "host:port to listen on")

func main() {
	flag.Parse()

	mux := http.NewServeMux()

	api := NewAPI()
	mux.Handle("/api/", http.StripPrefix("/api", api))

	srv := &http.Server{
		Addr:    *bindFlag,
		Handler: mux,
	}

	log.Printf("listening on %s", *bindFlag)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
