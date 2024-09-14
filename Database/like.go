package Database

import (
	"database/sql"
	"talknet/structs" // Adjust import path as needed
	"time"
)

// CreateLike inserts a new like into the database.
func CreateLike(db *sql.DB, userID int, postID *int, commentID *int) error {
	value := 1
	_, err := db.Exec("INSERT INTO likes_dislikes (user_id, post_id, comment_id, created_at, like_dislike) VALUES (?, ?, ?, ?, ?)",
		userID, postID, commentID, time.Now(), value)
	return err
}

// GetLikesByPostID retrieves likes for a post by its ID.
func GetReactionsByPostID(db *sql.DB, postID int) ([]structs.Like, []structs.Dislike, error) {
	// Query for likes (like_dislike = 1)
	likeRows, err := db.Query("SELECT id, user_id, post_id, comment_id, created_at FROM Likes_Dislikes WHERE post_id = ? AND like_dislike = 1", postID)
	if err != nil {
		return nil, nil, err
	}
	defer likeRows.Close()

	// Query for dislikes (like_dislike = 0)
	dislikeRows, err := db.Query("SELECT id, user_id, post_id, comment_id, created_at FROM Likes_Dislikes WHERE post_id = ? AND like_dislike = 0", postID)
	if err != nil {
		return nil, nil, err
	}
	defer dislikeRows.Close()

	var likes []structs.Like
	var dislikes []structs.Dislike

	// Scan likes data
	for likeRows.Next() {
		var like structs.Like
		err := likeRows.Scan(&like.ID, &like.UserID, &like.PostID, &like.CommentID, &like.CreatedAt)
		if err != nil {
			return nil, nil, err
		}
		likes = append(likes, like)
	}

	// Scan dislikes data
	for dislikeRows.Next() {
		var dislike structs.Dislike
		err := dislikeRows.Scan(&dislike.ID, &dislike.UserID, &dislike.PostID, &dislike.CommentID, &dislike.CreatedAt)
		if err != nil {
			return nil, nil, err
		}
		dislikes = append(dislikes, dislike)
	}

	return likes, dislikes, nil
}

// Other like-related functions (e.g., DeleteLike) go here.
func RemoveLikeDislike(db *sql.DB, userID int, postID int) error {
	_, err := db.Exec("DELETE FROM Likes_Dislikes WHERE user_id = ? AND post_id = ?", userID, postID)
	return err
}

func CreateDislike(db *sql.DB, userID int, postID *int, commentID *int) error {
	value := 0
	_, err := db.Exec("INSERT INTO likes_dislikes (user_id, post_id, comment_id, created_at, like_dislike) VALUES (?, ?, ?, ?, ?)",
		userID, postID, commentID, time.Now(), value)
	return err
}

func GetLikeDislikeCounts(db *sql.DB, postID int) (int, int, error) {
	var likeCount, dislikeCount int
	err := db.QueryRow("SELECT COUNT(*) FROM Likes_Dislikes WHERE post_id = ? AND like_dislike = 1 ", postID).Scan(&likeCount)
	if err != nil {
		return 0, 0, err
	}
	err = db.QueryRow("SELECT COUNT(*) FROM Likes_Dislikes WHERE post_id = ? AND like_dislike = 0 ", postID).Scan(&dislikeCount)
	if err != nil {
		return 0, 0, err
	}
	return likeCount, dislikeCount, nil
}

func CheckReactionExists(db *sql.DB, postID int, userID int) (int, error) {
	var value bool
	err := db.QueryRow("SELECT like_dislike FROM Likes_Dislikes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, nil // No reaction found
		}
		return -1, err
	}
	if value {
		return 1, nil // User has liked
	}
	return 0, nil // User has disliked
}
