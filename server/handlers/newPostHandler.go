package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
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
			return
		}
	}
	if r.Method == "POST" {
		id, isLoggedIn := sessions.GetSessionUserID(r)
		if !isLoggedIn {
			http.Error(w, "You Are Not Logged In", http.StatusBadRequest)
			return
		}
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
		// Retrieve form values
		title := r.FormValue("title")
		content := r.FormValue("content")
		categories := r.Form["category[]"] // Get the selected categories

		// Check for errors
		if id == -1 || title == "" || content == "" || len(categories) == 0 {
			http.Error(w, "Error: All fields must be filled and at least one category selected", http.StatusBadRequest)
			return
		}

		transaction, err := db.Begin()
		if err != nil {
			http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
			return
		}

		// Insert into Posts table
		res, err := transaction.Exec("INSERT INTO Posts (user_id, title, content) VALUES (?, ?, ?)", id, title, content)
		if err != nil {
			transaction.Rollback()
			http.Error(w, "Failed to insert post", http.StatusInternalServerError)
			return
		}

		// Get the last inserted post ID
		postID, err := res.LastInsertId()
		if err != nil {
			transaction.Rollback()
			http.Error(w, "Failed to get post ID", http.StatusInternalServerError)
			return
		}

		// Insert each selected category into Post_Categories
		for _, categoryIDStr := range categories {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				// Handle error
				http.Error(w, "Invalid category ID", http.StatusBadRequest)
				return
			}

			// Use categoryID in your SQL query
			_, err = transaction.Exec("INSERT INTO Post_Categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
			if err != nil {
				transaction.Rollback()
				http.Error(w, "Failed to insert post categories", http.StatusInternalServerError)
				return
			}
		}

		// Commit the transaction
		if err := transaction.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		// Successfully inserted post and categories
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}
