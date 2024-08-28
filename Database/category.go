package Database

import (
	"database/sql"
	"talknet/structs" // Adjust import path as needed
	"time"
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

func GetAllGategories(db *sql.DB) ([]structs.Category, error) {
	rows, err := db.Query("SELECT id,name FROM categories")
	if err != nil {
		return nil, err
	}
	var categories []structs.Category
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.ID,&category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	defer rows.Close()
	return categories, nil
}
