package Database

import (
    "database/sql"
    "talknet/structs" // Adjust import path as needed
)

// CreatePostCategory inserts a new post-category association into the database.
func CreatePostCategory(db *sql.DB, postID, categoryID int) error {
    _, err := db.Exec("INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)", postID, categoryID)
    return err
}

// GetCategoriesByPostID retrieves categories for a post by its ID.
func GetCategoriesByPostID(db *sql.DB, postID int) ([]structs.Category, error) {
    rows, err := db.Query(`
        SELECT c.id, c.name, c.created_at
        FROM categories c
        JOIN post_categories pc ON c.id = pc.category_id
        WHERE pc.post_id = ?`, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []structs.Category
    for rows.Next() {
        var category structs.Category
        err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt)
        if err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }
    return categories, nil
}

// Other post-category-related functions (e.g., DeletePostCategory) go here.
