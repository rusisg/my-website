package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"personal-website-template/internal/handlers"
	"time"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully waits for existing connections to finish - e.g., 15s or 1m")
	flag.Parse()

	mux := http.NewServeMux()
	handlers.RegisterRoutes(mux)

	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Static files: Images (including CV file)
	imgFileServer := http.FileServer(http.Dir("./assets/img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", imgFileServer))

	srv := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	go func() {
		log.Printf("Starting server at %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server gracefully stopped")
}
