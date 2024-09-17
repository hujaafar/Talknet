package Database

import (
	"database/sql"
	"talknet/structs"
	"time"
)

// CreateUser inserts a new user into the database.
func CreateUser(db *sql.DB, username, email, password string) error {
	_, err := db.Exec("INSERT INTO users (username, email, password, created_at) VALUES (?, ?, ?, ?)",
		username, email, password, time.Now())
	return err
}

// GetUserByUsername retrieves a user by username.

func GetUserByUsername(db *sql.DB, username string) (structs.User, error) {
	row := db.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE username = ?", username)

	var user structs.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

// function to validate username
func IsValidUsername(db *sql.DB, username string) bool {
	row := db.QueryRow("SELECT username FROM users WHERE username = ?", username)
	var user structs.User
	err := row.Scan(&user.Username)
	if err == sql.ErrNoRows {
		return true
	} else if err != nil {
		return false
	}
	return false
}

func GetUserByID(db *sql.DB, id int) (structs.User, error) {
	row := db.QueryRow("SELECT id, username, email, password, created_at FROM users WHERE id = ?", id)

	var user structs.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}
