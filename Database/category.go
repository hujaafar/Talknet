package Database

import (
    "database/sql"
    "time"
    "talknet/structs" // Adjust import path as needed
)

// CreateCategory inserts a new category into the database.
func CreateCategory(db *sql.DB, name string) error {
    _, err := db.Exec("INSERT INTO categories (name, created_at) VALUES (?, ?)", name, time.Now())
    return err
}

// GetCategoryByID retrieves a category by its ID.
func GetCategoryByID(db *sql.DB, id int) (structs.Category, error) {
    row := db.QueryRow("SELECT id, name, created_at FROM categories WHERE id = ?", id)
    var category structs.Category
    err := row.Scan(&category.ID, &category.Name, &category.CreatedAt)
    return category, err
}

// Other category-related functions (e.g., UpdateCategory, DeleteCategory) go here.
