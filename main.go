package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"talknet/server/handlers"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Open a connection to the database
	dbPath := "./talknet.db"       // Path to your SQLite database file
	sqlFilePath := "./talknet.sql" // Path to your SQL file

	var database *sql.DB

	// Check if the database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {

		// Create a new database
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		database = db

		// Read the SQL file
		sqlData, err := ioutil.ReadFile(sqlFilePath)
		if err != nil {
			log.Fatalf("Error reading SQL file: %v", err)
		}

		// Execute the SQL commands from the file
		_, err = database.Exec(string(sqlData))
		if err != nil {
			log.Fatalf("Error executing SQL commands: %v", err)
		}

	} else if err != nil {
		log.Fatalf("Error checking database file: %v", err)
	} else {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		database = db
	}

	// Ensure database is closed when main function exits
	defer database.Close()

	// Setup static file server
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Setup handlers and pass the database connection
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
	http.HandleFunc("/like_dislike", func(w http.ResponseWriter, r *http.Request) {
		handlers.LikeDislikeHandler(database, w, r)
	})
	http.HandleFunc("/post-details", func(w http.ResponseWriter, r *http.Request) {
		handlers.PostDetailsHandler(database, w, r)
	})
	http.HandleFunc("/add_comment", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddCommentHandler(database, w, r)
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		handlers.ProfileHandler(database, w, r)
	})
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		handlers.RenderErrorPage(w, "Error Message", http.StatusInternalServerError)
	})

	// Logout handler
	http.HandleFunc("/logout", handlers.LogoutHandler)

	// Start the server
	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
