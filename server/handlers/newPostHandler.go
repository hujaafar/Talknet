package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"talknet/Database"
	"talknet/server/sessions"
	"talknet/structs"
)

type Categories struct {
	AllCategories []structs.Category
}

func NewPostHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		http.NotFound(w, r)
		return
	}
	if r.Method == "GET" {
		allCategories, err := Database.GetAllGategories(db)
		if err != nil {
			log.Printf("Failed to get all categories: %v", err)
			http.Error(w, "Failed to load categories", http.StatusInternalServerError)
			return
		}
		categories := Categories{
			AllCategories: allCategories,
		}
		err = templates.ExecuteTemplate(w, "new-post.html", categories)
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
	if r.Method == "POST" {
		id,isLoggedIn:= sessions.GetSessionUserID(r)
		if !isLoggedIn {
			http.Error(w, "You Are Not Logged In", http.StatusBadRequest)
		}
		title:= r.Form.Get("title")
		content := r.Form.Get("content")
		if id==-1||title==""||content==""{
			http.Error(w, "error", http.StatusBadRequest)
		}
		
	}
}
