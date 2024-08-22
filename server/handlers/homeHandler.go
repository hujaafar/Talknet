package handlers

import (
	"html/template"
	"log"
	"net/http"

	"talknet/server/sessions"
)

type TemplateData struct {
	IsLoggedIn bool
}

var templates = template.Must(template.ParseGlob("static/*.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, isLoggedIn := sessions.GetSessionUserID(r)

	data := TemplateData{
		IsLoggedIn: isLoggedIn,
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
