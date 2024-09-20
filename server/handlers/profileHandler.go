package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"talknet/Database"
	"talknet/server/sessions"
)

func ProfileHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var myPostDataList []PostData
	var likedPostDataList []PostData

	if r.URL.Path != "/profile" {
		RenderErrorPage(w, "Not Found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		RenderErrorPage(w, "Not Found", http.StatusNotFound)
		return
	}

	userID, isLoggedIn := sessions.GetSessionUserID(r)

	// Check if he requests his profile or someone else's profile
	var profileID int
	if r.URL.Query().Get("id") == "" {
		profileID = userID
	} else {
		postID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("Failed to parse profile ID: %v", err)
			RenderErrorPage(w, "Invalid profile ID", http.StatusBadRequest)
			return
		}
		profileID, err = Database.GetUserIdByPostID(db, postID)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			RenderErrorPage(w, "You are not logged in", http.StatusInternalServerError)
			return
		}
	}

	username, err := Database.GetUsername(db, profileID)
	if err != nil {
		log.Printf("Failed to get username: %v", err)
		RenderErrorPage(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	isHisProfile := profileID == userID

	// Fetch My Posts
	posts, err := Database.GetPostByUserID(db, profileID)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		RenderErrorPage(w, "Failed to load posts", http.StatusInternalServerError)
		return
	}

	// Fetch Liked Posts
	likedPosts, err := Database.GetLikedPosts(db, profileID)
	if err != nil {
		log.Printf("Failed to get liked posts: %v", err)
		RenderErrorPage(w, "Failed to load liked posts", http.StatusInternalServerError)
		return
	}

	// Process My Posts
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
			reaction, err = Database.CheckPostReactionExists(db, post.ID, userID)
			if err != nil {
				log.Printf("Failed to check reaction: %v", err)
				continue
			}
		}

		myPostDataList = append(myPostDataList, PostData{
			ID:             post.ID,
			Username:       user.Username,
			Title:          post.Title,
			Content:        post.Content,
			CreatedAt:      timeAgo(post.CreatedAt),
			PostCategories: postCategories,
			LikeCount:      likeCount,
			DislikeCount:   dislikeCount,
			CommentCount:   len(comments),
			Reaction:       reaction,
		})
	}

	// Process Liked Posts
	for _, post := range likedPosts {
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
			reaction, err = Database.CheckPostReactionExists(db, post.ID, userID)
			if err != nil {
				log.Printf("Failed to check reaction: %v", err)
				continue
			}
		}

		likedPostDataList = append(likedPostDataList, PostData{
			ID:             post.ID,
			Username:       user.Username,
			Title:          post.Title,
			Content:        post.Content,
			CreatedAt:      timeAgo(post.CreatedAt),
			PostCategories: postCategories,
			LikeCount:      likeCount,
			DislikeCount:   dislikeCount,
			CommentCount:   len(comments),
			Reaction:       reaction,
		})
	}

	// Combine both lists into a single data structure
	data := struct {
		MyPosts      []PostData
		LikedPosts   []PostData
		IsHisProfile bool
		Username     string
	}{
		MyPosts:      myPostDataList,
		LikedPosts:   likedPostDataList,
		IsHisProfile: isHisProfile,
		Username:     username,
	}

	// Render the profile template with both my posts and liked posts
	err = templates.ExecuteTemplate(w, "Profile.html", data)
	if err != nil {
		log.Printf("Failed to render template: %v", err)
		RenderErrorPage(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
