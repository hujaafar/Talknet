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
func GetCategoryNamesByPostID(db *sql.DB, postID int) ([]structs.Category, error) {
	query := `
		SELECT c.name 
		FROM Post_Categories pc
		JOIN Categories c ON pc.category_id = c.id
		WHERE pc.post_id = ?
	`

	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []structs.Category
	for rows.Next() {
		var category structs.Category
		if err := rows.Scan(&category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Log the retrieved categories for debugging

	return categories, nil
}
