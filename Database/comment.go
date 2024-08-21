package Database

import (
    "database/sql"
    "time"
    "talknet/structs" // Adjust import path as needed
)

// CreateComment inserts a new comment into the database.
func CreateComment(db *sql.DB, postID, userID int, content string) error {
    _, err := db.Exec("INSERT INTO comments (post_id, user_id, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        postID, userID, content, time.Now(), time.Now())
    return err
}

// GetCommentByID retrieves a comment by its ID.
func GetCommentByID(db *sql.DB, id int) (structs.Comment, error) {
    row := db.QueryRow("SELECT id, post_id, user_id, content, created_at, updated_at FROM comments WHERE id = ?", id)
    var comment structs.Comment
    err := row.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
    return comment, err
}

// Other comment-related functions (e.g., UpdateComment, DeleteComment) go here.
