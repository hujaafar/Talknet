package server

import (
	"database/sql"
	 "talknet/Database"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, username, email, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
	if  !Database.IsValidUsername(db, username) {
		return err
	}
    return Database.CreateUser(db, username, email, string(hashedPassword))
}
