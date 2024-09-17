package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"talknet/Database"
	"talknet/structs"

)

var postDetailTemplate = template.Must(template.ParseFiles("static/pages/post-details.html"))

type CommentWithUser struct {
	structs.Comment
	Username  string
	CreatedAt string
}

func PostDetailsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL query
	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Fetch the post by ID
	post, err := Database.GetPostByID(db, postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Fetch the user who created the post
	user, err := Database.GetUserByID(db, post.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Fetch comments for the post
	comments, err := Database.GetCommentsByPostID(db, postID)
	if err != nil {
		log.Printf("Failed to get comments: %v", err)
		http.Error(w, "Failed to load comments", http.StatusInternalServerError)
		return
	}

	// Prepare comments with usernames
	var commentsWithUser []CommentWithUser
	for _, comment := range comments {
		commentUser, err := Database.GetUserByID(db, comment.UserID)
		if err != nil {
			log.Printf("Failed to get user for comment: %v", err)
			continue
		}
		commentsWithUser = append(commentsWithUser, CommentWithUser{
			Comment:   comment,
			Username:  commentUser.Username,
			CreatedAt: timeAgo(comment.CreatedAt),
		})
	}

	// Render the post details template
	err = postDetailTemplate.Execute(w, struct {
		Post     structs.Post
		Username string
		Comments []CommentWithUser
	}{
		Post:     post,
		Username: user.Username,
		Comments: commentsWithUser,
	})
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

