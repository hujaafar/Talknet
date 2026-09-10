package forum

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"talknet/static"
	"time"
)

type App struct {
	attempts authAttempts
	db       *sql.DB
	secure   bool
	views    *template.Template
	handler  http.Handler
}
type csrfKey struct{}

func csrfValue(r *http.Request) string {
	value, _ := r.Context().Value(csrfKey{}).(string)
	return value
}

type Page struct {
	Notice                                         string
	PageNumber                                     int
	PreviousURL, NextURL                           string
	Title, Error, CSRF, Query, Category, Sort, Tab string
	User                                           *User
	Profile                                        *User
	Categories                                     []Category
	Posts                                          []Post
	Post                                           *Post
	Comments                                       []Comment
	Values                                         map[string]string
	TotalPosts, TotalMembers                       int
	Status                                         int
}

func New(db *sql.DB, secure bool) *App {
	funcs := template.FuncMap{
		"initial": func(s string) string {
			r := []rune(s)
			if len(r) == 0 {
				return "?"
			}
			return strings.ToUpper(string(r[0]))
		},
		"date": func(t time.Time) string { return t.Format("Jan 2, 2006") },
		"ago": func(t time.Time) string {
			d := time.Since(t)
			if d < time.Hour {
				return "Recently"
			}
			if d < 24*time.Hour {
				return fmt.Sprintf("%dh ago", int(d.Hours()))
			}
			if d < 7*24*time.Hour {
				return fmt.Sprintf("%dd ago", int(d.Hours()/24))
			}
			return t.Format("Jan 2, 2006")
		},
		"excerpt": func(s string) string {
			r := []rune(s)
			if len(r) > 170 {
				return string(r[:170]) + "…"
			}
			return s
		},
	}
	a := &App{db: db, secure: secure, attempts: authAttempts{entries: make(map[string]attempt)}, views: template.Must(template.New("talknet").Funcs(funcs).ParseFS(static.Files, "pages/*.html"))}
	mux := http.NewServeMux()
	// Only public assets are served; embedded templates are never exposed as source.
	for _, dir := range []string{"styles", "js", "images"} {
		mux.Handle("GET /static/"+dir+"/", http.StripPrefix("/static/", http.FileServer(http.FS(static.Files))))
	}
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.PingContext(r.Context()); err != nil {
			http.Error(w, "unavailable", 503)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, "ok")
	})
	mux.HandleFunc("GET /{$}", a.home)
	mux.HandleFunc("GET /post-details", a.detail)
	mux.HandleFunc("GET /profile", a.profile)
	mux.HandleFunc("GET /login", a.login)
	mux.HandleFunc("POST /login", a.login)
	mux.HandleFunc("GET /register", a.register)
	mux.HandleFunc("POST /register", a.register)
	mux.HandleFunc("GET /post", a.compose)
	mux.HandleFunc("POST /post", a.compose)
	mux.HandleFunc("POST /add_comment", a.comment)
	mux.HandleFunc("POST /like_dislike", a.react)
	mux.HandleFunc("POST /logout", a.logout)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { a.fail(w, r, 404, "That page couldn’t be found.") })
	a.handler = a.security(http.NewCrossOriginProtection().Handler(mux))
	return a
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) { a.handler.ServeHTTP(w, r) }
func randomToken() string {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func hashToken(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func (a *App) cookie(w http.ResponseWriter, name, value string, maxAge int) {
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Path: "/", HttpOnly: true, Secure: a.secure, SameSite: http.SameSiteLaxMode, MaxAge: maxAge})
}

func (a *App) security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'")
		if a.secure {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000")
		}
		if strings.HasPrefix(r.URL.Path, "/static/") || r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-store")
		c, err := r.Cookie("csrf_token")
		if err != nil || len(c.Value) != 64 {
			c = &http.Cookie{Name: "csrf_token", Value: randomToken()}
			a.cookie(w, c.Name, c.Value, 86400)
		}
		r = r.WithContext(context.WithValue(r.Context(), csrfKey{}, c.Value))
		r.Body = http.MaxBytesReader(w, r.Body, 65536)
		if r.Method != "GET" && r.Method != "HEAD" {
			token := r.Header.Get("X-CSRF-Token")
			if token == "" {
				if err := r.ParseForm(); err != nil {
					a.fail(w, r, 400, "The request is too large or could not be read.")
					return
				}
				token = r.PostForm.Get("csrf_token")
			}
			if subtle.ConstantTimeCompare([]byte(token), []byte(c.Value)) != 1 {
				a.fail(w, r, 403, "Your form expired. Refresh the page and try again.")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) user(r *http.Request) *User {
	c, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}
	var u User
	err = a.db.QueryRow(`SELECT u.id,u.username,u.email,u.created_at FROM Users u JOIN Sessions s ON u.id=s.user_id WHERE s.session_token=? AND s.expires_at>?`, hashToken(c.Value), time.Now().UTC()).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil
	}
	return &u
}
func (a *App) page(r *http.Request, title string) Page {
	return Page{Title: title, User: a.user(r), CSRF: csrfValue(r), Values: map[string]string{}}
}
func uid(u *User) int {
	if u == nil {
		return 0
	}
	return u.ID
}
func (a *App) render(w http.ResponseWriter, name string, p Page, status int) {
	var buf bytes.Buffer
	if err := a.views.ExecuteTemplate(&buf, name, p); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "Unable to display this page", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	_, _ = buf.WriteTo(w)
}
func (a *App) fail(w http.ResponseWriter, r *http.Request, status int, message string) {
	p := Page{Title: "Something went wrong", Status: status, Error: message, CSRF: csrfValue(r)}
	a.render(w, "error.html", p, status)
}
func (a *App) internal(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s: %v", r.URL.Path, err)
	a.fail(w, r, 500, "We couldn’t load this right now. Please try again.")
}
func (a *App) requireUser(w http.ResponseWriter, r *http.Request) *User {
	u := a.user(r)
	if u == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	return u
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	p := a.page(r, "A place for curious minds")
	p.Query = strings.TrimSpace(r.URL.Query().Get("q"))
	p.Category = r.URL.Query().Get("category")
	p.Sort = r.URL.Query().Get("sort")
	if p.Sort != "popular" {
		p.Sort = "latest"
	}
	if len(p.Query) > 200 {
		a.fail(w, r, 400, "Keep your search under 200 characters.")
		return
	}
	var err error
	p.Categories, err = a.categories()
	if err != nil {
		a.internal(w, r, err)
		return
	}
	p.PageNumber = pageNumber(r)
	p.Posts, err = a.posts(postQuery{Viewer: uid(p.User), Category: p.Category, Search: p.Query, Sort: p.Sort, Page: p.PageNumber})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	if err = a.db.QueryRow("SELECT (SELECT COUNT(*) FROM Posts),(SELECT COUNT(*) FROM Users)").Scan(&p.TotalPosts, &p.TotalMembers); err != nil {
		a.internal(w, r, err)
		return
	}
	setPagination(r, &p)
	a.render(w, "index.html", p, 200)
}

func (a *App) detail(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("post_id"))
	if err != nil || id < 1 {
		a.fail(w, r, 400, "Choose a valid discussion.")
		return
	}
	p := a.page(r, "Discussion")
	posts, err := a.posts(postQuery{Viewer: uid(p.User), ID: id})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	for i := range posts {
		if posts[i].ID == id {
			p.Post = &posts[i]
			break
		}
	}
	if p.Post == nil {
		a.fail(w, r, 404, "This discussion doesn’t exist.")
		return
	}
	p.Title = p.Post.Title
	rows, err := a.db.Query(`SELECT c.id,c.user_id,c.content,u.username,c.created_at,
 (SELECT COUNT(*) FROM Likes_Dislikes WHERE comment_id=c.id AND like_dislike=1),
 (SELECT COUNT(*) FROM Likes_Dislikes WHERE comment_id=c.id AND like_dislike=0),
 COALESCE((SELECT like_dislike FROM Likes_Dislikes WHERE comment_id=c.id AND user_id=? LIMIT 1),-1)
 FROM Comments c JOIN Users u ON c.user_id=u.id WHERE c.post_id=? ORDER BY c.id`, uid(p.User), id)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	for rows.Next() {
		var c Comment
		if err = rows.Scan(&c.ID, &c.UserID, &c.Content, &c.Username, &c.CreatedAt, &c.Likes, &c.Dislikes, &c.Reaction); err != nil {
			break
		}
		p.Comments = append(p.Comments, c)
	}
	if err == nil {
		err = rows.Err()
	}
	rows.Close()
	if err != nil {
		a.internal(w, r, err)
		return
	}
	a.render(w, "post-details.html", p, 200)
}

func (a *App) profile(w http.ResponseWriter, r *http.Request) {
	p := a.page(r, "Member profile")
	id := uid(p.User)
	if raw := r.URL.Query().Get("user"); raw != "" {
		var err error
		id, err = strconv.Atoi(raw)
		if err != nil {
			a.fail(w, r, 400, "Choose a valid member.")
			return
		}
	}
	// Preserve old profile links, which used a post ID rather than a user ID.
	if raw := r.URL.Query().Get("id"); raw != "" {
		if err := a.db.QueryRow("SELECT user_id FROM Posts WHERE id=?", raw).Scan(&id); err != nil {
			a.fail(w, r, 404, "Member not found.")
			return
		}
	}
	if id == 0 {
		http.Redirect(w, r, "/login", 303)
		return
	}
	var u User
	err := a.db.QueryRow("SELECT id,username,email,created_at FROM Users WHERE id=?", id).Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
	if err == sql.ErrNoRows {
		a.fail(w, r, 404, "Member not found.")
		return
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	p.Profile = &u
	p.Title = u.Username
	p.Tab = r.URL.Query().Get("tab")
	// Liked discussions are personal; other members only see authored posts.
	liked := p.Tab == "liked" && id == uid(p.User)
	if !liked {
		p.Tab = "posts"
	}
	p.PageNumber = pageNumber(r)
	p.Posts, err = a.posts(postQuery{Viewer: uid(p.User), Author: id, Liked: liked, Page: p.PageNumber})
	if err != nil {
		a.internal(w, r, err)
		return
	}
	setPagination(r, &p)
	a.render(w, "Profile.html", p, 200)
}

func pageNumber(r *http.Request) int {
	n, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || n < 1 || n > 100000 {
		return 1
	}
	return n
}
func setPagination(r *http.Request, p *Page) {
	link := func(n int) string {
		q := r.URL.Query()
		q.Set("page", strconv.Itoa(n))
		fragment := ""
		if r.URL.Path == "/" {
			fragment = "#discussions"
		}
		return r.URL.Path + "?" + q.Encode() + fragment
	}
	if p.PageNumber > 1 {
		p.PreviousURL = link(p.PageNumber - 1)
	}
	if len(p.Posts) > pageSize {
		p.NextURL = link(p.PageNumber + 1)
		p.Posts = p.Posts[:pageSize]
	}
}
