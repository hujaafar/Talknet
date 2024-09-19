package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"talknet/server"
	"talknet/server/sessions"
)

type Data struct {
	ErrorMsg string
}

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
			loginData := Data{ErrorMsg: "Invalid Username or Password"}
			err := templates.ExecuteTemplate(w, "login.html", loginData)
			if err != nil {
				log.Printf("Failed to render template: %v", err)
				RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
			}
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

			data := Data{ErrorMsg: "Email is already taken"}

			err := templates.ExecuteTemplate(w, "SignUp.html", data)
			if err != nil {
				log.Printf("Failed to render template: %v", err)
				RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
