package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"talknet/server/handlers"

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
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(database, w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginHandler(database, w, r)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.RegisterHandler(database, w, r)
	})
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		handlers.NewPostHandler(database, w, r)
	})
	http.HandleFunc("/logout", handlers.LogoutHandler)
	fmt.Println("Server running at http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
