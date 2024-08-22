package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type TemplateData struct {
	IsLoggedIn bool
}

var templates = template.Must(template.ParseGlob("static/*.html"))

func checkUserSession(r *http.Request) bool {
	// Example: Check if a session cookie exists and is valid
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		return false
	}

	// Here you can add logic to validate the session token, e.g., checking it against a database or an in-memory store.
	// For simplicity, let's assume the session token is valid if it exists.
	return true
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	isLoggedIn := checkUserSession(r)

	data := TemplateData{
		IsLoggedIn: isLoggedIn,
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
