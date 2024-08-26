package Database

import (
	"database/sql"
	"talknet/structs" // Adjust import path as needed
	"time"
)

// CreateComment inserts a new comment into the database.
func CreateComment(db *sql.DB, postID, userID int, content string) error {
	_, err := db.Exec("INSERT INTO comments (post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?)",
		postID, userID, content, time.Now(), time.Now())
	return err
}

// GetCommentByID retrieves a comment by its ID.

func GetCommentByID(db *sql.DB, postID int) ([]structs.Comment, error) {
	rows, err := db.Query("SELECT id, post_id, user_id, content, created_at FROM comments WHERE post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []structs.Comment
	for rows.Next() {
		var comment structs.Comment
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
