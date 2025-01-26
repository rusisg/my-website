package handlers

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"personal-website-template/internal/json"
	"personal-website-template/internal/lib/token"
	"text/template"
	"time"
)

type Login struct {
	SessionToken string
	CSRFToken    string
}

var logged = map[string]Login{}

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
	// Check if the session token cookie exists
	cookie, err := r.Cookie("token")
	if err == nil {
		// Validate the session token
		for _, session := range logged {
			if session.SessionToken == cookie.Value {
				// If valid, redirect to the admin page
				http.Redirect(w, r, "/note/admin/new", http.StatusSeeOther)
				log.Println("User already logged in. Redirecting to /note/admin/new")
				return
			}
		}
	}

	switch r.Method {
	case http.MethodGet:
		// Render the login form
		files := []string{
			"./assets/html/login.page.gohtml",
			"./assets/html/base.layout.gohtml",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Println("Error parsing templates:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return

	case http.MethodPost:
		// Process login form submission
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form submission", http.StatusBadRequest)
			return
		}

		username := r.FormValue("admin_username")
		password := r.FormValue("admin_password")

		// Load .env variables
		err = godotenv.Load(".env")
		if err != nil {
			log.Println("Error loading .env file:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		envUser := os.Getenv("ADMIN_USER")
		envPass := os.Getenv("ADMIN_PASS")

		// Validate credentials
		if username != envUser || password != envPass {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		// Generate tokens
		sessionToken := token.GenerateToken(32)
		csrfToken := token.GenerateToken(32)

		// Set cookies
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    sessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false, // Accessible to client-side
			Secure:   true,
		})

		// Store session in memory
		logged[username] = Login{
			SessionToken: sessionToken,
			CSRFToken:    csrfToken,
		}

		// Redirect to the admin page
		http.Redirect(w, r, "/note/admin/new", http.StatusSeeOther)
		log.Println("User logged in:", username)
		return

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func NoteAdminNew(w http.ResponseWriter, r *http.Request) {
	// Handle POST request
	if r.Method == http.MethodGet {
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
		log.Println(r.Method, "note admin new page")
	} else if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
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
	} else if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the session token is valid
		for _, session := range logged {
			if session.SessionToken == cookie.Value {
				next(w, r)
				return
			}
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/note", Note)
	mux.HandleFunc("/note/admin", NoteAdmin)

	// Wrap the /note/admin/new route with the Authenticate middleware
	mux.Handle("/note/admin/new", AuthMiddleware(NoteAdminNew))
}
