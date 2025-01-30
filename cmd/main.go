package main

import (
	"log"
	"net/http"
	"os"
	"personal-website-template/internal/handlers"
	"personal-website-template/internal/lib/graceful_shutdown"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Falling back to system environment variables.")
	}

	mux := http.NewServeMux()

	// Register application routes
	handlers.RegisterRoutes(mux)

	// Serve static files: CSS, JS
	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Serve image files (including CV)
	imgFileServer := http.FileServer(http.Dir("./assets/img/"))
	mux.Handle("/img/", http.StripPrefix("/img/", imgFileServer))

	// Get port and host from environment variables or use defaults
	port := os.Getenv("ALWAYSDATA_HTTPD_PORT")
	if port == "" {
		log.Println("ALWAYSDATA_HTTPD_PORT not set. Falling back to default port 8100.")
		port = "8100"
	}

	host := os.Getenv("ALWAYSDATA_HTTPD_IP")
	if host == "" {
		log.Println("ALWAYSDATA_HTTPD_IP not set. Falling back to default host 0.0.0.0.")
		host = "0.0.0.0"
	}

	// Wrap IPv6 addresses in square brackets
	if host == "::" || host == "0:0:0:0:0:0:0:1" {
		host = "[" + host + "]"
	}

	addr := host + ":" + port

	// Create HTTP server
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Handle graceful shutdown
	err = graceful_shutdown.GracefulShutdown(srv)
	if err != nil {
		log.Fatalf("Failed to gracefully shut down: %s", err)
	}
}
