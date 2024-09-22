package server

import (
	"database/sql"
	"errors"
	"regexp"
	"talknet/Database"
	"golang.org/x/crypto/bcrypt"
)

// Regular expressions
var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\.]{3,}$`)                   // At least 3 characters, only alphanumeric, underscores, and dots
	passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`)   // At least 8 characters, must include upper, lower case letters, and digits
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`) // Valid email format
)

func RegisterUser(db *sql.DB, username, email, password string) error {
	// Validate username
	if !usernameRegex.MatchString(username) {
		return errors.New("Invalid username: must be at least 3 characters long and can only contain letters, numbers, underscores, and dots")
	}

	// Validate email
	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format")
	}

	// Validate password
	if !passwordRegex.MatchString(password) {
		return errors.New("Invalid password: must be at least 8 characters long, contain an uppercase letter, a lowercase letter, and a digit")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Check if the username is valid in the database
	if !Database.IsValidUsername(db, username) {
		return errors.New("Username is not valid")
	}

	// Create the user in the database
	return Database.CreateUser(db, username, email, string(hashedPassword))
}
