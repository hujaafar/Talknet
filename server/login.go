package server

import (
	"database/sql"
	"errors"
	"talknet/Database"
	"talknet/structs"
	"golang.org/x/crypto/bcrypt"
)

// Use the same regular expressions as for registration
func LoginUser(db *sql.DB, username, password string) (structs.User, error) {
	// Validate username
	if !usernameRegex.MatchString(username) {
		return structs.User{}, errors.New("Invalid username")
	}

	// Validate password
	if !passwordRegex.MatchString(password) {
		return structs.User{}, errors.New("Invalid password")
	}

	// Fetch user by username
	user, err := Database.GetUserByUsername(db, username)
	if err != nil {
		return user, err
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil
}
