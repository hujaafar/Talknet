package Database

import (
    "database/sql"
    "time"
    "talknet/structs" // Adjust import path as needed
)

// CreatePost inserts a new post into the database.
func CreatePost(db *sql.DB, userID int, title, content string) error {
    _, err := db.Exec("INSERT INTO posts (user_id, title, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
        userID, title, content, time.Now(), time.Now())
    return err
}

// GetPostByID retrieves a post by its ID.
func GetPostByID(db *sql.DB, id int) (structs.Post, error) {
    row := db.QueryRow("SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE id = ?", id)
    var post structs.Post
    err := row.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
    return post, err
}

// Other post-related functions (e.g., UpdatePost, DeletePost) go here.
