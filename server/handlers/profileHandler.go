package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"talknet/Database"
	"talknet/server/sessions"
)

func ProfileHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var postDataList []PostData

	if r.URL.Path != "/profile" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	userID, isLoggedIn := sessions.GetSessionUserID(r)
	if !isLoggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	posts, err := Database.GetPostByUserID(db, userID)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	for _, post := range posts {
		user, err := Database.GetUserByID(db, post.UserID)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			continue
		}

		postCategories, err := Database.GetCategoryNamesByPostID(db, post.ID)
		if err != nil {
			log.Printf("Failed to get categories: %v", err)
			continue
		}

		likes, dislikes, err := Database.GetReactionsByPostID(db, post.ID)
		if err != nil {
			log.Printf("Failed to get likes: %v", err)
			continue
		}
		likeCount := len(likes)
		dislikeCount := len(dislikes)
		comments, err := Database.GetCommentsByPostID(db, post.ID)
		if err != nil {
			log.Printf("Failed to get comments: %v", err)
			continue
		}
		reaction := -1
		if isLoggedIn {
			reaction, err = Database.CheckReactionExists(db, post.ID, userID)
			if err != nil {
				log.Printf("Failed to check reaction: %v", err)
				continue
			}
		}

		postDataList = append(postDataList, PostData{
			ID:             post.ID,
			Username:       user.Username,
			Title:          post.Title,
			Content:        post.Content,
			CreatedAt:      timeAgo(post.CreatedAt), // Use relative time format
			PostCategories: postCategories,
			LikeCount:      likeCount,
			DislikeCount:   dislikeCount,
			CommentCount:   len(comments),
			Reaction:       reaction,
		})
	}

	err = templates.ExecuteTemplate(w, "Profile.html", postDataList)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
