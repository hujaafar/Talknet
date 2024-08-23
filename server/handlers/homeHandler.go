package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"talknet/Database"
	"talknet/server/sessions"
	"talknet/structs"
	"time"
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

type PostData struct {
	ID           int
	Username     string
	Title        string
	Content      string
	CreatedAt    string
	Categories   []structs.Category
	LikeCount    int
	DislikeCount int
	CommentCount int
}

func IndexHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	posts, err := Database.GetAllPosts(db)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	var postData []PostData

	for _, post := range posts {
		user, err := Database.GetUserByUsername(db, string(post.UserID))
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			continue
		}

		categories, err := Database.GetCategoriesByPostID(db, post.ID)
		if err != nil {
			log.Printf("Failed to get categories: %v", err)
			continue
		}

		likes, err := Database.GetLikesByPostID(db, post.ID)
		if err != nil {
			log.Printf("Failed to get likes: %v", err)
			continue
		}
		likeCount := len(likes)

		comments, err := Database.GetCommentByID(db, post.ID) 
		if err != nil {
			log.Printf("Failed to get comments: %v", err)
			continue
		}

		postData = append(postData, PostData{
			ID:           post.ID,
			Username:     user.Username,
			Title:        post.Title,
			Content:      post.Content,
			CreatedAt:    post.CreatedAt.Format(time.RFC822),
			Categories:   categories,
			LikeCount:    likeCount,
			CommentCount: len(comments), 
		})
	}

	err = templates.ExecuteTemplate(w, "index.html", postData)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
