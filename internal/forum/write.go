package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (a *App) compose(w http.ResponseWriter, r *http.Request) {
	u := a.requireUser(w, r)
	if u == nil {
		return
	}
	p := a.page(r, "Start a conversation")
	p.FormAction = r.URL.Path
	p.Editing = r.URL.Path == "/post/edit"
	if p.Editing {
		id, err := strconv.Atoi(r.FormValue("post_id"))
		if err != nil || id < 1 {
			a.fail(w, r, 400, "Choose a valid discussion.")
			return
		}
		posts, err := a.posts(postQuery{Viewer: u.ID, ID: id})
		if err != nil {
			a.internal(w, r, err)
			return
		}
		if len(posts) == 0 {
			a.fail(w, r, 404, "This discussion doesn’t exist.")
			return
		}
		p.Post = &posts[0]
		if p.Post.UserID != u.ID {
			a.fail(w, r, 403, "Only the author can edit this discussion.")
			return
		}
		p.Title = "Edit your discussion"
		p.Values["title"] = p.Post.Title
		p.Values["content"] = p.Post.Content
		p.Values["revision"] = strconv.Itoa(p.Post.Revision)
		for _, c := range p.Post.Categories {
			p.Values[fmt.Sprintf("category_%d", c.ID)] = "selected"
		}
	}
	var err error
	p.Categories, err = a.categories()
	if err != nil {
		a.internal(w, r, err)
		return
	}
	if r.Method == "GET" {
		a.render(w, "new-post.html", p, 200)
		return
	}
	title := strings.TrimSpace(r.FormValue("title"))
	content := strings.TrimSpace(r.FormValue("content"))
	// Submitted values replace defaults, including topics the author unchecked.
	p.Values = map[string]string{"revision": r.FormValue("revision")}
	p.Values["title"] = title
	p.Values["content"] = content
	ids := []int{}
	seen := map[int]bool{}
	valid := map[int]bool{}
	for _, c := range p.Categories {
		valid[c.ID] = true
	}
	for _, raw := range r.PostForm["category[]"] {
		id, err := strconv.Atoi(raw)
		if err != nil || !valid[id] {
			p.Error = "Choose an available topic."
			break
		}
		if !seen[id] {
			ids = append(ids, id)
			seen[id] = true
			p.Values[fmt.Sprintf("category_%d", id)] = "selected"
		}
	}
	if utf8.RuneCountInString(title) < 5 || utf8.RuneCountInString(title) > 120 {
		p.Error = "Give your discussion a title between 5 and 120 characters."
	} else if utf8.RuneCountInString(content) < 10 || utf8.RuneCountInString(content) > 5000 {
		p.Error = "Write between 10 and 5,000 characters."
	} else if len(ids) < 1 || len(ids) > 3 {
		p.Error = "Choose between one and three topics."
	}
	if p.Error != "" {
		a.render(w, "new-post.html", p, 400)
		return
	}
	tx, err := a.db.Begin()
	if err != nil {
		a.internal(w, r, err)
		return
	}
	defer tx.Rollback()
	var id int64
	if p.Editing {
		id = int64(p.Post.ID)
		var owner, revision int
		err = tx.QueryRow("SELECT user_id,COALESCE((SELECT revision FROM Post_Revisions WHERE post_id=Posts.id),1) FROM Posts WHERE id=?", id).Scan(&owner, &revision)
		if err != nil {
			a.internal(w, r, err)
			return
		}
		if owner != u.ID {
			a.fail(w, r, 403, "Only the author can edit this discussion.")
			return
		}
		expected, parseErr := strconv.Atoi(r.FormValue("revision"))
		if parseErr != nil || expected != revision {
			p.Error = "This discussion changed in another tab. Your text is kept below. Open the current discussion before applying your changes."
			a.render(w, "new-post.html", p, 409)
			return
		}
		_, err = tx.Exec("UPDATE Posts SET title=?,content=? WHERE id=? AND user_id=?", title, content, id, u.ID)
		if err == nil {
			_, err = tx.Exec("DELETE FROM Post_Categories WHERE post_id=?", id)
		}
		if err == nil {
			_, err = tx.Exec("INSERT INTO Post_Revisions(post_id,revision) VALUES(?,2) ON CONFLICT(post_id) DO UPDATE SET revision=revision+1", id)
		}
	} else {
		var result sql.Result
		result, err = tx.Exec("INSERT INTO Posts(user_id,title,content) VALUES(?,?,?)", u.ID, title, content)
		if err == nil {
			id, err = result.LastInsertId()
		}
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	for _, c := range ids {
		if _, err = tx.Exec("INSERT INTO Post_Categories(post_id,category_id) VALUES(?,?)", id, c); err != nil {
			a.internal(w, r, err)
			return
		}
	}
	if err = tx.Commit(); err != nil {
		a.internal(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post-details?post_id=%d", id), 303)
}

func (a *App) comment(w http.ResponseWriter, r *http.Request) {
	u := a.requireUser(w, r)
	if u == nil {
		return
	}
	id, err := strconv.Atoi(r.FormValue("post_id"))
	content := strings.TrimSpace(r.FormValue("content"))
	if err != nil || id < 1 {
		a.fail(w, r, 400, "Choose a valid discussion.")
		return
	}
	if utf8.RuneCountInString(content) < 1 || utf8.RuneCountInString(content) > 2000 {
		a.fail(w, r, 400, "Write a reply between 1 and 2,000 characters.")
		return
	}
	var exists int
	if err = a.db.QueryRow("SELECT id FROM Posts WHERE id=?", id).Scan(&exists); err == sql.ErrNoRows {
		a.fail(w, r, 404, "This discussion doesn’t exist.")
		return
	} else if err != nil {
		a.internal(w, r, err)
		return
	}
	if _, err = a.db.Exec("INSERT INTO Comments(post_id,user_id,content) VALUES(?,?,?)", id, u.ID, content); err != nil {
		a.internal(w, r, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post-details?post_id=%d#replies", id), 303)
}

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
func (a *App) react(w http.ResponseWriter, r *http.Request) {
	u := a.user(r)
	if u == nil {
		jsonResponse(w, 401, map[string]string{"error": "Sign in to react to a discussion."})
		return
	}
	var in struct {
		PostID int    `json:"postId"`
		Action string `json:"action"`
		Type   string `json:"type"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&in); err != nil {
		jsonResponse(w, 400, map[string]string{"error": "Invalid reaction."})
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		jsonResponse(w, 400, map[string]string{"error": "Invalid reaction."})
		return
	}
	if in.PostID < 1 || (in.Action != "like" && in.Action != "dislike") || (in.Type != "post" && in.Type != "comment") {
		jsonResponse(w, 400, map[string]string{"error": "Invalid reaction."})
		return
	}
	// SQL identifiers only come from this closed allowlist; user values stay bound.
	table, column := "Posts", "post_id"
	if in.Type == "comment" {
		table, column = "Comments", "comment_id"
	}
	tx, err := a.db.Begin()
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "Couldn’t save your reaction."})
		return
	}
	defer tx.Rollback()
	var id int
	err = tx.QueryRow("SELECT id FROM "+table+" WHERE id=?", in.PostID).Scan(&id)
	if err == sql.ErrNoRows {
		jsonResponse(w, 404, map[string]string{"error": "This conversation no longer exists."})
		return
	}
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "Couldn’t save your reaction."})
		return
	}
	previous := -1
	err = tx.QueryRow("SELECT CAST(like_dislike AS INTEGER) FROM Likes_Dislikes WHERE user_id=? AND "+column+"=? LIMIT 1", u.ID, id).Scan(&previous)
	if err != nil && err != sql.ErrNoRows {
		jsonResponse(w, 500, map[string]string{"error": "Couldn’t save your reaction."})
		return
	}
	desired := 1
	if in.Action == "dislike" {
		desired = 0
	}
	_, err = tx.Exec("DELETE FROM Likes_Dislikes WHERE user_id=? AND "+column+"=?", u.ID, id)
	if err == nil && previous != desired {
		_, err = tx.Exec("INSERT INTO Likes_Dislikes(user_id,"+column+",like_dislike) VALUES(?,?,?)", u.ID, id, desired)
	} else if previous == desired {
		desired = -1
	}
	var likes, dislikes int
	if err == nil {
		err = tx.QueryRow("SELECT COALESCE(SUM(like_dislike=1),0),COALESCE(SUM(like_dislike=0),0) FROM Likes_Dislikes WHERE "+column+"=?", id).Scan(&likes, &dislikes)
	}
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		jsonResponse(w, 500, map[string]string{"error": "Couldn’t save your reaction."})
		return
	}
	jsonResponse(w, 200, map[string]int{"likeCount": likes, "dislikeCount": dislikes, "reaction": desired})
}
