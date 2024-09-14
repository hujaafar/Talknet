package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"talknet/Database"
	"talknet/server/sessions"
)

func LikeDislikeHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		PostID int    `json:"postId"`
		Action string `json:"action"` // "like" or "dislike"
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userID, isLoggedIn := sessions.GetSessionUserID(r)
	if !isLoggedIn {
		log.Println("User not logged in")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Remove any existing like/dislike by this user on this post
	err = Database.RemoveLikeDislike(db, userID, requestData.PostID)
	if err != nil {
		log.Println("Error removing existing like/dislike:", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// Add new like/dislike
	log.Println("Action:", requestData.Action)
	if requestData.Action == "like" {
		err = Database.CreateLike(db, userID, &requestData.PostID, nil)
	} else if requestData.Action == "dislike" {
		err = Database.CreateDislike(db, userID, &requestData.PostID, nil)
	}

	if err != nil {
		log.Println("Error creating like/dislike:", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// Get updated like/dislike counts
	likeCount, dislikeCount, err := Database.GetLikeDislikeCounts(db, requestData.PostID)
	if err != nil {
		log.Println("Error getting like/dislike counts:", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// Send the updated counts back to the client
	responseData := map[string]interface{}{
		"likeCount":    likeCount,
		"dislikeCount": dislikeCount,
		"action":       requestData.Action,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
