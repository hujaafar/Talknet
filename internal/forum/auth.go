package forum

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"net"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,24}$`)

type attempt struct {
	count int
	until time.Time
}

type authAttempts struct {
	sync.Mutex
	entries map[string]attempt
}

func (a *App) allowAuth(r *http.Request) bool {
	attempts := &a.attempts
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	attempts.Lock()
	defer attempts.Unlock()
	now := time.Now()
	for k, v := range attempts.entries {
		if now.After(v.until) {
			delete(attempts.entries, k)
		}
	}
	v := attempts.entries[ip]
	if v.count == 0 {
		v.until = now.Add(time.Minute)
	}
	if v.count >= 10 || len(attempts.entries) >= 10000 {
		return false
	}
	v.count++
	attempts.entries[ip] = v
	return true
}

func (a *App) login(w http.ResponseWriter, r *http.Request) {
	p := a.page(r, "Welcome back")
	if p.User != nil {
		http.Redirect(w, r, "/", 303)
		return
	}
	if r.Method == "GET" {
		if r.URL.Query().Get("registered") == "1" {
			p.Notice = "Your account is ready. Sign in to join the conversation."
		}
		a.render(w, "login.html", p, 200)
		return
	}
	if !a.allowAuth(r) {
		w.Header().Set("Retry-After", "60")
		a.fail(w, r, 429, "Too many attempts. Wait a minute and try again.")
		return
	}
	name := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	p.Values["username"] = name
	var id int
	var hash string
	err := a.db.QueryRow("SELECT id,password FROM Users WHERE username=?", name).Scan(&id, &hash)
	if err != nil && err != sql.ErrNoRows {
		a.internal(w, r, err)
		return
	}
	// Match the expensive comparison on missing accounts to reduce username probing.
	if err == sql.ErrNoRows {
		hash = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy"
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil || id == 0 {
		p.Error = "That username and password don’t match."
		a.render(w, "login.html", p, 401)
		return
	}
	token := randomToken()
	tx, err := a.db.Begin()
	if err != nil {
		a.internal(w, r, err)
		return
	}
	defer tx.Rollback()
	if _, err = tx.Exec("DELETE FROM Sessions WHERE user_id=? OR expires_at<=?", id, time.Now().UTC()); err == nil {
		_, err = tx.Exec("INSERT INTO Sessions(user_id,session_token,expires_at) VALUES(?,?,?)", id, hashToken(token), time.Now().UTC().Add(24*time.Hour))
	}
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		a.internal(w, r, err)
		return
	}
	a.cookie(w, "session_id", token, 86400)
	http.Redirect(w, r, "/", 303)
}

func (a *App) register(w http.ResponseWriter, r *http.Request) {
	p := a.page(r, "Find your people")
	if p.User != nil {
		http.Redirect(w, r, "/", 303)
		return
	}
	if r.Method == "GET" {
		a.render(w, "SignUp.html", p, 200)
		return
	}
	if !a.allowAuth(r) {
		w.Header().Set("Retry-After", "60")
		a.fail(w, r, 429, "Too many attempts. Wait a minute and try again.")
		return
	}
	name := strings.TrimSpace(r.FormValue("username"))
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")
	p.Values["username"] = name
	p.Values["email"] = email
	address, err := mail.ParseAddress(email)
	switch {
	case !usernamePattern.MatchString(name):
		p.Error = "Use 3–24 letters, numbers, or underscores for your username."
	case err != nil || address.Address != email || len(email) > 254:
		p.Error = "Enter a valid email address."
	case utf8.RuneCountInString(password) < 12:
		p.Error = "Use a password with at least 12 characters."
	case len(password) > 72:
		p.Error = "This password is too long. Try a shorter passphrase."
	}
	if p.Error != "" {
		a.render(w, "SignUp.html", p, 400)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.internal(w, r, err)
		return
	}
	_, err = a.db.Exec("INSERT INTO Users(username,email,password) VALUES(?,?,?)", name, email, string(hash))
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code == sqlite3.ErrConstraint {
			p.Error = "That username or email is already registered."
			a.render(w, "SignUp.html", p, 409)
			return
		}
		a.internal(w, r, err)
		return
	}
	http.Redirect(w, r, "/login?registered=1", 303)
}

func (a *App) logout(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie("session_id"); err == nil {
		if _, err := a.db.Exec("DELETE FROM Sessions WHERE session_token=?", hashToken(c.Value)); err != nil {
			a.internal(w, r, err)
			return
		}
	}
	a.cookie(w, "session_id", "", -1)
	http.Redirect(w, r, "/", 303)
}
