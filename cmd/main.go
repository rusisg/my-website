package main

import (
	"log"
	"net/http"
	"personal-website-template/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Home)
	mux.HandleFunc("/note", handlers.Note)

	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	log.Printf("Starting server at %s", ":4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
