package handlers

import (
	"log"
	"net/http"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)

	files := []string{
		"./assets/html/home.page.gohtml",
		"./assets/html/base.layout.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	log.Println("home page")
}

func Note(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/note" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)

	files := []string{
		"./assets/html/note.page.gohtml",
		"./assets/html/base.layout.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	log.Println("note page")
}

func NoteAdmin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/note/admin" {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
	files := []string{
		"./assets/html/note_admin.page.gohtml",
		"./assets/html/base.layout.gohtml",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	log.Println("note admin page")
}
