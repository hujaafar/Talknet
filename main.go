package main

import (
    "database/sql"
    "log"
    _ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
    // Open a connection to the database
    database, err := sql.Open("sqlite3", "./talknet.db")
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    // Test the connection
    err = database.Ping()
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    } else {
        log.Println("Connected to the database successfully!")
    }
}
