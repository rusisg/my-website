package web

import "net/http"

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Home page"))
}

func Blog(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Blog page"))
}
