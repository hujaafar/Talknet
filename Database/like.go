package Database

import (
    "database/sql"
    "time"
    "talknet/structs" // Adjust import path as needed
)

// CreateLike inserts a new like into the database.
func CreateLike(db *sql.DB, userID int, postID *int, commentID *int) error {
    _, err := db.Exec("INSERT INTO likes (user_id, post_id, comment_id, created_at) VALUES (?, ?, ?, ?)",
        userID, postID, commentID, time.Now())
    return err
}

// GetLikesByPostID retrieves likes for a post by its ID.
func GetLikesByPostID(db *sql.DB, postID int) ([]structs.Like, error) {
    rows, err := db.Query("SELECT id, user_id, post_id, comment_id, created_at FROM likes WHERE post_id = ?", postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var likes []structs.Like
    for rows.Next() {
        var like structs.Like
        err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.CommentID, &like.CreatedAt)
        if err != nil {
            return nil, err
        }
        likes = append(likes, like)
    }
    return likes, nil
}

// Other like-related functions (e.g., DeleteLike) go here.
