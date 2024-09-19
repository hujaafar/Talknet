package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"talknet/server"
	"talknet/server/sessions"
)


// Login Handler
func LoginHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		// Process the login form
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := server.LoginUser(db, username, password)
		if err != nil {
			RenderErrorPage(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		sessions.CreateSession(w, user.ID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Register Handler
func RegisterHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "SignUp.html", nil)
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		// Process the registration form
		r.ParseForm()
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")
		err := server.RegisterUser(db, username, email, password)
		if err != nil {
			RenderErrorPage(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
