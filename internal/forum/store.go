package forum

import (
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"time"
)

//go:embed schema.sql
var schema string

func Open(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on&_busy_timeout=5000&_journal_mode=WAL")
	if err != nil {
		return nil, err
	}
	// A single connection serializes SQLite writes. Close rows before nested queries.
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("initialize database: %w", err)
	}
	return db, nil
}

type Category struct {
	ID    int
	Name  string
	Count int
}
type User struct {
	ID              int
	Username, Email string
	CreatedAt       time.Time
}
type Post struct {
	Saved                               bool
	Revision                            int
	ID, UserID                          int
	Title, Content, Username            string
	CreatedAt                           time.Time
	Categories                          []Category
	Likes, Dislikes, Comments, Reaction int
}
type Comment struct {
	ID, UserID                int
	Content, Username         string
	CreatedAt                 time.Time
	Likes, Dislikes, Reaction int
}

func (a *App) categories() ([]Category, error) {
	rows, err := a.db.Query("SELECT c.id,c.name,COUNT(pc.post_id) FROM Categories c LEFT JOIN Post_Categories pc ON pc.category_id=c.id GROUP BY c.id ORDER BY c.name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Count); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

type postQuery struct {
	Viewer, Author, ID, Page int
	Category, Search, Sort   string
	Liked                    bool
	Saved                    bool
}

const pageSize = 20

func (a *App) posts(filter postQuery) ([]Post, error) {
	query := `SELECT p.id,p.user_id,p.title,p.content,u.username,p.created_at,
 (SELECT COUNT(*) FROM Likes_Dislikes WHERE post_id=p.id AND like_dislike=1),
 (SELECT COUNT(*) FROM Likes_Dislikes WHERE post_id=p.id AND like_dislike=0),
 (SELECT COUNT(*) FROM Comments WHERE post_id=p.id),
 COALESCE((SELECT like_dislike FROM Likes_Dislikes WHERE post_id=p.id AND user_id=? LIMIT 1),-1),
 EXISTS(SELECT 1 FROM Bookmarks WHERE post_id=p.id AND user_id=?),
 COALESCE((SELECT revision FROM Post_Revisions WHERE post_id=p.id),1)
 FROM Posts p JOIN Users u ON u.id=p.user_id WHERE 1=1`
	args := []any{filter.Viewer, filter.Viewer}
	if filter.ID > 0 {
		query += " AND p.id=?"
		args = append(args, filter.ID)
	}
	if filter.Category != "" {
		query += " AND EXISTS(SELECT 1 FROM Post_Categories pc JOIN Categories c ON c.id=pc.category_id WHERE pc.post_id=p.id AND c.name=?)"
		args = append(args, filter.Category)
	}
	if filter.Search != "" {
		query += " AND (instr(lower(p.title),lower(?))>0 OR instr(lower(p.content),lower(?))>0)"
		args = append(args, filter.Search, filter.Search)
	}
	if filter.Author > 0 {
		if filter.Saved {
			query += " AND EXISTS(SELECT 1 FROM Bookmarks WHERE post_id=p.id AND user_id=?)"
		} else if filter.Liked {
			query += " AND EXISTS(SELECT 1 FROM Likes_Dislikes WHERE post_id=p.id AND user_id=? AND like_dislike=1)"
		} else {
			query += " AND p.user_id=?"
		}
		args = append(args, filter.Author)
	}
	if filter.Saved && filter.Author > 0 {
		// Preserve save order even when several saves share the same second.
		query += " ORDER BY (SELECT created_at FROM Bookmarks WHERE post_id=p.id AND user_id=?) DESC,(SELECT rowid FROM Bookmarks WHERE post_id=p.id AND user_id=?) DESC"
		args = append(args, filter.Author, filter.Author)
	} else if filter.Sort == "popular" {
		query += " ORDER BY 7 DESC,9 DESC,p.id DESC"
	} else {
		query += " ORDER BY p.id DESC"
	}
	if filter.Page < 1 {
		filter.Page = 1
	}
	// Fetch one extra row to discover the next page without a second count query.
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize+1, (filter.Page-1)*pageSize)
	rows, err := a.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	var out []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Title, &p.Content, &p.Username, &p.CreatedAt, &p.Likes, &p.Dislikes, &p.Comments, &p.Reaction, &p.Saved, &p.Revision); err != nil {
			rows.Close()
			return nil, err
		}
		out = append(out, p)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return nil, err
	}
	for i := range out {
		rows, err := a.db.Query("SELECT c.id,c.name FROM Categories c JOIN Post_Categories pc ON pc.category_id=c.id WHERE pc.post_id=? ORDER BY c.name", out[i].ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var c Category
			if err := rows.Scan(&c.ID, &c.Name); err != nil {
				rows.Close()
				return nil, err
			}
			out[i].Categories = append(out[i].Categories, c)
		}
		err = rows.Err()
		rows.Close()
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
