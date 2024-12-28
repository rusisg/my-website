package middleware

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			log.Printf("Received credentials - Username: %s, Password: %s", username, password)

			err := godotenv.Load(".env")
			if err != nil {
				log.Printf("Failed to load .env: %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			envUser := os.Getenv("ADMIN_USER")
			envPass := os.Getenv("ADMIN_PASS")
			log.Printf("Expected credentials - Username: %s, Password: %s", envUser, envPass)

			if username != envUser || password != envPass {
				http.Error(w, "Who are you? :)", http.StatusForbidden)
				return
			}
		})
}
