package forum

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Search exposes only the public metadata needed by the command palette.
func (a *App) quickSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if utf8.RuneCountInString(q) > 200 {
		jsonResponse(w, 400, map[string]string{"error": "Keep searches under 200 characters."})
		return
	}
	type result struct {
		ID     int    `json:"id"`
		Title  string `json:"title"`
		Author string `json:"author"`
		Topic  string `json:"topic"`
	}
	out := []result{}
	if len([]rune(q)) >= 2 {
		posts, err := a.posts(postQuery{Search: q})
		if err != nil {
			jsonResponse(w, 500, map[string]string{"error": "Search is unavailable. Try again."})
			return
		}
		for i, p := range posts {
			if i == 6 {
				break
			}
			topic := "Discussion"
			if len(p.Categories) > 0 {
				topic = p.Categories[0].Name
			}
			out = append(out, result{p.ID, p.Title, p.Username, topic})
		}
	}
	jsonResponse(w, 200, map[string]any{"results": out})
}

// Explicit save/remove actions are safe to repeat after a network retry.
func (a *App) bookmark(w http.ResponseWriter, r *http.Request) {
	json := strings.Contains(r.Header.Get("Accept"), "application/json")
	fail := func(status int, message string) {
		if json {
			jsonResponse(w, status, map[string]string{"error": message})
		} else {
			a.fail(w, r, status, message)
		}
	}
	u := a.user(r)
	if u == nil {
		if json {
			fail(401, "Sign in to save discussions.")
		} else {
			http.Redirect(w, r, "/login", 303)
		}
		return
	}
	id, err := strconv.Atoi(r.FormValue("post_id"))
	action := r.FormValue("action")
	if err != nil || id < 1 || (action != "save" && action != "remove") {
		fail(400, "Choose a valid discussion and save action.")
		return
	}
	tx, err := a.db.Begin()
	if err != nil {
		fail(500, "Couldn’t update your saved discussions.")
		return
	}
	defer tx.Rollback()
	var exists int
	err = tx.QueryRow("SELECT id FROM Posts WHERE id=?", id).Scan(&exists)
	if err == sql.ErrNoRows {
		fail(404, "This discussion no longer exists.")
		return
	}
	if err != nil {
		fail(500, "Couldn’t update your saved discussions.")
		return
	}
	if action == "save" {
		_, err = tx.Exec("INSERT OR IGNORE INTO Bookmarks(user_id,post_id) VALUES(?,?)", u.ID, id)
	} else {
		_, err = tx.Exec("DELETE FROM Bookmarks WHERE user_id=? AND post_id=?", u.ID, id)
	}
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		fail(500, "Couldn’t update your saved discussions.")
		return
	}
	if json {
		jsonResponse(w, 200, map[string]bool{"saved": action == "save"})
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post-details?post_id=%d", id), 303)
}
