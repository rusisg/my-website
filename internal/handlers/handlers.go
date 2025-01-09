package handlers

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"personal-website-template/internal/json"
	"personal-website-template/internal/middleware"
	"text/template"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

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

// TODO: It must get datas and show it from database

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
	if r.Method == http.MethodGet {
		// Render the login form
		files := []string{
			"./assets/html/login.page.gohtml",
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
	} else if r.Method == http.MethodPost {
		// Process login form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}

		username := r.FormValue("admin_username")
		password := r.FormValue("admin_password")

		err = godotenv.Load(".env")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		envUser := os.Getenv("ADMIN_USER")
		envPass := os.Getenv("ADMIN_PASS")

		if username == envUser && password == envPass {
			// Credentials are valid, redirect to /note/admin/new
			http.Redirect(w, r, "/note/admin/new", http.StatusFound)
		} else {
			// Invalid credentials
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		}
	}
}

func NoteAdminNew(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}

		// Retrieve data from the form
		title := r.FormValue("title")            // ID of the title input
		content := r.FormValue("text-paragraph") // ID of the textarea

		// Create a new note
		newNote := json.NoteStruct{
			Title:   title,
			Content: content,
		}

		// Read existing notes
		notes, err := json.ReadNotes()
		if err != nil {
			log.Printf("Failed to read notes: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		// Append the new note and save to the JSON file
		notes = append(notes, newNote)
		err = json.WriteNotes(notes)
		if err != nil {
			log.Printf("Failed to save note: %v", err)
			http.Error(w, "Internal Server Error", 500)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Note saved successfully"))
	}

}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/note", Note)
	mux.HandleFunc("/note/admin", NoteAdmin)

	// Wrap the /note/admin/new route with the Authenticate middleware
	mux.Handle("/note/admin/new", middleware.Authenticate(http.HandlerFunc(NoteAdminNew)))
}
