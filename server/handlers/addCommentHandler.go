package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"talknet/Database"
	"talknet/server/sessions"
)

func AddCommentHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Check if the user is logged in
	userID, isLoggedIn := sessions.GetSessionUserID(r)
	if !isLoggedIn {
		http.Error(w, "You must be logged in to comment", http.StatusUnauthorized)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get form values
	content := r.FormValue("content")
	postIDStr := r.FormValue("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Validate content
	if content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	// Save the comment to the database
	err = Database.CreateComment(db, postID, userID, content)
	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	// Redirect back to the post details page
	http.Redirect(w, r, "/post-details?post_id="+postIDStr, http.StatusSeeOther)
}