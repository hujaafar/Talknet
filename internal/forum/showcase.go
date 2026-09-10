package forum

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//go:embed showcase.json
var showcaseJSON []byte

type DemoSummary struct {
	Users, Posts, Replies, Likes int
	AlreadyAdded                 bool
}

// SeedShowcase is an explicit additive import. Existing records are never
// adopted, updated, or deleted; an identity collision rolls back the whole pack.
func SeedShowcase(db *sql.DB) (DemoSummary, error) {
	var pack struct {
		Authors []string `json:"authors"`
		Posts   []struct {
			Title, Content, Category string
			Author                   int
			Replies                  []string
		} `json:"posts"`
	}
	if err := json.Unmarshal(showcaseJSON, &pack); err != nil {
		return DemoSummary{}, err
	}
	tx, err := db.Begin()
	if err != nil {
		return DemoSummary{}, err
	}
	defer tx.Rollback()
	var imported int
	err = tx.QueryRow("SELECT 1 FROM Demo_Seeds WHERE name=?", "showcase-v1").Scan(&imported)
	if err == nil {
		return DemoSummary{AlreadyAdded: true}, nil
	}
	if err != sql.ErrNoRows {
		return DemoSummary{}, err
	}
	// Reserve the pack in the same transaction as its data.
	if _, err = tx.Exec("INSERT INTO Demo_Seeds(name) VALUES(?)", "showcase-v1"); err != nil {
		return DemoSummary{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(randomToken()), bcrypt.DefaultCost)
	if err != nil {
		return DemoSummary{}, err
	}
	now := time.Now().UTC()
	authors := make([]int64, len(pack.Authors))
	summary := DemoSummary{}
	for i, name := range pack.Authors {
		result, err := tx.Exec("INSERT INTO Users(username,email,password,created_at) VALUES(?,?,?,?)", name, strings.ToLower(name)+"@example.invalid", string(hash), now.Add(-time.Duration(30-i*2)*24*time.Hour))
		if err != nil {
			return DemoSummary{}, fmt.Errorf("add fictional author %s: %w", name, err)
		}
		authors[i], err = result.LastInsertId()
		if err != nil {
			return DemoSummary{}, err
		}
		summary.Users++
	}
	for i, post := range pack.Posts {
		if post.Author < 0 || post.Author >= len(authors) {
			return DemoSummary{}, fmt.Errorf("invalid author in showcase post %d", i)
		}
		var category int
		if err = tx.QueryRow("SELECT id FROM Categories WHERE name=?", post.Category).Scan(&category); err != nil {
			return DemoSummary{}, fmt.Errorf("showcase topic %s: %w", post.Category, err)
		}
		created := now.Add(-time.Hour - time.Duration(len(pack.Posts)-1-i)*75*time.Minute)
		result, err := tx.Exec("INSERT INTO Posts(user_id,title,content,created_at) VALUES(?,?,?,?)", authors[post.Author], post.Title, post.Content, created)
		if err != nil {
			return DemoSummary{}, err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return DemoSummary{}, err
		}
		if _, err = tx.Exec("INSERT INTO Post_Categories(post_id,category_id) VALUES(?,?)", id, category); err != nil {
			return DemoSummary{}, err
		}
		summary.Posts++
		for j, reply := range post.Replies {
			author := authors[(post.Author+1+j*2)%len(authors)]
			if _, err = tx.Exec("INSERT INTO Comments(post_id,user_id,content,created_at) VALUES(?,?,?,?)", id, author, reply, created.Add(time.Duration(15+j*20)*time.Minute)); err != nil {
				return DemoSummary{}, err
			}
			summary.Replies++
		}
		for j := 1; j <= 2+i%4; j++ {
			if _, err = tx.Exec("INSERT INTO Likes_Dislikes(user_id,post_id,like_dislike,created_at) VALUES(?,?,1,?)", authors[(post.Author+j)%len(authors)], id, created.Add(40*time.Minute)); err != nil {
				return DemoSummary{}, err
			}
			summary.Likes++
		}
	}
	if err = tx.Commit(); err != nil {
		return DemoSummary{}, err
	}
	return summary, nil
}
