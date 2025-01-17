package main

import (
	"log"
	"net/http"
	"personal-website-template/internal/handlers"
	"personal-website-template/internal/lib"
)

func main() {
	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Static files: Images (including CV file)
	imgFileServer := http.FileServer(http.Dir("./assets/img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", imgFileServer))

	srv := &http.Server{
		Addr:    ":5000",
		Handler: mux,
	}

	err := lib.GracefulShutdown(srv)
	if err != nil {
		log.Fatalf("Failed to graceful shutdown: %s", err)
	}
}
