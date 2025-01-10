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

	log.Println(r.Method, "home page")
}

func Note(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/note" {
		http.NotFound(w, r)
		return
	}

	notes, err := json.ReadNotes()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	files := []string{
		"./assets/html/note.page.gohtml",
		"./assets/html/base.layout.gohtml",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Parse: Internal Server Error", 500)
		return
	}

	data := struct {
		Notes []json.NoteStruct
	}{
		Notes: notes,
	}

	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	log.Println(r.Method, "note page")
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
	log.Println(r.Method, "note admin page")
}

func NoteAdminNew(w http.ResponseWriter, r *http.Request) {
	// Handle POST request
	if r.Method == http.MethodPost {
		// Render the note creation form
		files := []string{
			"./assets/html/new_note.page.gohtml",
			"./assets/html/base.layout.gohtml",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Printf("Error parsing templates: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Execute the template
		err = ts.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Parse the form data
		err = r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}

		// Retrieve data from the form
		title := r.FormValue("title")
		content := r.FormValue("text-paragraph")

		// Validate inputs
		if title == "" || content == "" {
			http.Error(w, "Title and content are required", http.StatusBadRequest)
			return
		}

		// Read existing notes
		notes, err := json.ReadNotes()
		if err != nil {
			log.Printf("Failed to read notes: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var newID int64
		if len(notes) > 0 {
			newID = notes[0].ID + 1 // Increment based on the latest (first) note
		} else {
			newID = 1 // Start from 1 if there are no notes
		}

		// Create a new note
		newNote := json.NoteStruct{
			ID:      newID,
			Title:   title,
			Content: content,
		}

		// Append the new note and save to the JSON file
		notes = append([]json.NoteStruct{newNote}, notes...)

		// Save the notes to the JSON file
		err = json.WriteNotes(notes)
		if err != nil {
			log.Printf("Failed to save note: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		log.Println(r.Method, "note admin new note page")

		// Redirect to the form after successful operation
		http.Redirect(w, r, "/note/admin/new", http.StatusSeeOther)
		return
	}

	// Handle unsupported HTTP methods (Method Not Allowed)
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/note", Note)
	mux.HandleFunc("/note/admin", NoteAdmin)

	// Wrap the /note/admin/new route with the Authenticate middleware
	mux.Handle("/note/admin/new", middleware.LoggingMiddleware(http.HandlerFunc(NoteAdminNew)))
}
