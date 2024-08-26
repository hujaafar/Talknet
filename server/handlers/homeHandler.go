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

type StaticPageData struct {
	IsLoggedIn    bool
	AllCategories []structs.Category
}

type PostData struct {
	ID             int
	Username       string
	Title          string
	Content        string
	CreatedAt      string
	PostCategories []structs.Category
	LikeCount      int
	DislikeCount   int
	CommentCount   int
}


var templates = template.Must(template.ParseGlob("static/*.html"))

func HomeHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, isLoggedIn := sessions.GetSessionUserID(r)

	// Fetch static data
	allCategories, err := Database.GetAllGategories(db)
	if err != nil {
		log.Printf("Failed to get all categories: %v", err)
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	staticData := StaticPageData{
		IsLoggedIn:    isLoggedIn,
		AllCategories: allCategories,
	}

	// Fetch dynamic post data
	posts, err := Database.GetAllPosts(db)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		http.Error(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	var postDataList []PostData

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

		postDataList = append(postDataList, PostData{
			ID:             post.ID,
			Username:       user.Username,
			Title:          post.Title,
			Content:        post.Content,
			CreatedAt:      post.CreatedAt.Format(time.RFC822),
			PostCategories: postCategories,
			LikeCount:      likeCount,
			CommentCount:   len(comments),
		})
	}

	// Render the template with both static and dynamic data
	err = templates.ExecuteTemplate(w, "index.html", struct {
		StaticData StaticPageData
		Posts      []PostData
	}{
		StaticData: staticData,
		Posts:      postDataList,
	})
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
