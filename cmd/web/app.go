package web

import (
	"log"
	"net/http"
)

const port = ":4000"

func App() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/blog", Blog)

	log.Printf("Starting server at %s", port)
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
