package server

import (
	"database/sql"
	"talknet/Database"
	"talknet/structs"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB, username, password string) (structs.User, error) {
    user, err := Database.GetUserByUsername(db, username)
    if err != nil {
        return user, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return user, err
    }

    return user, nil
}
