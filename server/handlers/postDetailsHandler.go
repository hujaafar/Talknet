package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"talknet/Database"
	"talknet/server/sessions"
	"talknet/structs"
)

var postDetailTemplate = template.Must(template.ParseFiles("static/pages/post-details.html"))

type CommentWithUser struct {
	structs.Comment
	Username     string
	CreatedAt    string
	LikeCount    int
	DislikeCount int
	CommentCount int
	Reaction     int
}

func PostDetailsHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Get the post ID from the URL query
	userSessionID, isLoggedIn := sessions.GetSessionUserID(r)

	postIDStr := r.URL.Query().Get("post_id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		RenderErrorPage(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Fetch the post by ID
	post, err := Database.GetPostByID(db, postID)
	if err != nil {
		RenderErrorPage(w, "Post not found", http.StatusNotFound)
		return
	}

	// Fetch the user who created the post
	user, err := Database.GetUserByID(db, post.UserID)
	if err != nil {
		RenderErrorPage(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Fetch comments for the post
	comments, err := Database.GetCommentsByPostID(db, postID)
	if err != nil {
		log.Printf("Failed to get comments: %v", err)
		RenderErrorPage(w, "Failed to load comments", http.StatusInternalServerError)
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

		likes, dislikes, err := Database.GetReactionsByCommentID(db, comment.ID)
		if err != nil {
			log.Printf("Failed to get likes: %v", err)
			continue
		}
		likeCount := len(likes)
		dislikeCount := len(dislikes)

		reaction := -1
		if isLoggedIn {
			reaction, err = Database.CheckReactionExists(db, comment.ID, userSessionID,"comment")
			if err != nil {
				log.Printf("Failed to check reaction: %v", err)
				continue
			}
		}

		commentsWithUser = append(commentsWithUser, CommentWithUser{
			Comment:      comment,
			Username:     commentUser.Username,
			CreatedAt:    timeAgo(comment.CreatedAt),
			LikeCount:    likeCount,
			DislikeCount: dislikeCount,
			Reaction:     reaction,
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
		RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
