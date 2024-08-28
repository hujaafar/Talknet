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
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		// Process the login form
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := server.LoginUser(db, username, password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		// Process the registration form
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		if err := server.RegisterUser(db, username, email, password); err != nil {
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
		http.ResponseWriter.Write(w, []byte("Registration successful!"))
		//http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
