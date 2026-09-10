package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

type testSite struct {
	db     *sql.DB
	server *httptest.Server
	client *http.Client
	base   *url.URL
}

func newSite(t *testing.T) *testSite {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(New(db, false))
	jar, _ := cookiejar.New(nil)
	base, _ := url.Parse(server.URL)
	s := &testSite{db, server, &http.Client{Jar: jar, Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}, base}
	t.Cleanup(func() { server.Close(); db.Close() })
	return s
}
func (s *testSite) csrf() string {
	for _, c := range s.client.Jar.Cookies(s.base) {
		if c.Name == "csrf_token" {
			return c.Value
		}
	}
	return ""
}
func (s *testSite) request(t *testing.T, method, path, body, contentType string, token bool) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequest(method, s.server.URL+path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if token {
		req.Header.Set("X-CSRF-Token", s.csrf())
	}
	resp, err := s.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	return resp, string(data)
}
func (s *testSite) form(t *testing.T, path string, values url.Values) (*http.Response, string) {
	t.Helper()
	values.Set("csrf_token", s.csrf())
	return s.request(t, "POST", path, values.Encode(), "application/x-www-form-urlencoded", false)
}
func checkStatus(t *testing.T, r *http.Response, want int, body string) {
	t.Helper()
	if r.StatusCode != want {
		t.Fatalf("status %d want %d: %s", r.StatusCode, want, body)
	}
}
func (s *testSite) signup(t *testing.T) {
	t.Helper()
	r, b := s.request(t, "GET", "/register", "", "", false)
	checkStatus(t, r, 200, b)
	r, b = s.form(t, "/register", url.Values{"username": {"ReaderOne"}, "email": {"reader@example.com"}, "password": {"A-long-passphrase-123"}})
	checkStatus(t, r, 303, b)
	r, b = s.form(t, "/login", url.Values{"username": {"ReaderOne"}, "password": {"A-long-passphrase-123"}})
	checkStatus(t, r, 303, b)
}

func TestCommunityFlow(t *testing.T) {
	s := newSite(t)
	for _, path := range []string{"/", "/login", "/register", "/healthz", "/static/styles/app.css", "/static/js/app.js", "/static/js/motion.mjs", "/static/js/theme.js", "/static/js/experience.mjs", "/static/styles/experience.css", "/static/images/conversation-art.webp", "/static/images/inter-latin.woff2", "/static/images/instrument-serif-italic.ttf"} {
		r, b := s.request(t, "GET", path, "", "", false)
		checkStatus(t, r, 200, b)
	}
	r, b := s.request(t, "GET", "/static/pages/index.html", "", "", false)
	checkStatus(t, r, 404, b)
	r, b = s.form(t, "/post", url.Values{"title": {"Guest post"}, "content": {"This should never be saved"}, "category[]": {"1"}})
	checkStatus(t, r, 303, b)
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM Posts").Scan(&count)
	if count != 0 {
		t.Fatal("guest created a post")
	}
	s.signup(t)
	r, b = s.form(t, "/post", url.Values{"title": {"A useful question"}, "content": {"How do you learn? <script>alert(1)</script>"}, "category[]": {"1", "1"}})
	checkStatus(t, r, 303, b)
	location := r.Header.Get("Location")
	r, b = s.request(t, "GET", location, "", "", false)
	checkStatus(t, r, 200, b)
	if strings.Contains(b, "<script>alert(1)</script>") || !strings.Contains(b, "&lt;script&gt;") {
		t.Fatal("post content was not escaped")
	}
	r, b = s.request(t, "GET", "/?q=useful&category=Technology&sort=popular", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, "A useful question") {
		t.Fatal("search/filter omitted matching post")
	}
	r, b = s.request(t, "GET", "/?q=nomatch", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, `class="empty-state"`) || strings.Contains(b, `class="post-card"`) {
		t.Fatal("missing empty state")
	}
	r, b = s.form(t, "/add_comment", url.Values{"post_id": {"1"}, "content": {"I learn by building small projects."}})
	checkStatus(t, r, 303, b)
	r, b = s.request(t, "GET", location, "", "", false)
	checkStatus(t, r, 200, b)
	for _, label := range []string{"Like reply", "Dislike reply"} {
		if !regexp.MustCompile(`aria-label="` + label + `"\s+aria-pressed="false"`).MatchString(b) {
			t.Fatalf("%s must render an exact false ARIA token before a reaction", label)
		}
	}
	for _, reaction := range []struct {
		action string
		want   int
	}{{"like", 1}, {"dislike", 0}, {"dislike", -1}} {
		r, b = s.request(t, "POST", "/like_dislike", `{"postId":1,"action":"`+reaction.action+`","type":"post"}`, "application/json", true)
		checkStatus(t, r, 200, b)
		var data struct {
			Reaction int `json:"reaction"`
		}
		if err := json.Unmarshal([]byte(b), &data); err != nil {
			t.Fatal(err)
		}
		if data.Reaction != reaction.want {
			t.Fatalf("reaction %d want %d", data.Reaction, reaction.want)
		}
	}
	r, b = s.request(t, "POST", "/like_dislike", `{"postId":1,"action":"like","type":"comment"}`, "application/json", true)
	checkStatus(t, r, 200, b)
	r, b = s.request(t, "GET", location, "", "", false)
	checkStatus(t, r, 200, b)
	if !regexp.MustCompile(`aria-label="Like reply"\s+aria-pressed="true"`).MatchString(b) || !regexp.MustCompile(`aria-label="Dislike reply"\s+aria-pressed="false"`).MatchString(b) {
		t.Fatal("reply buttons do not reflect the saved reaction after a page load")
	}
	r, b = s.request(t, "POST", "/like_dislike", `{"postId":1,"action":"like","type":"post"}`, "application/json", true)
	checkStatus(t, r, 200, b)
	r, b = s.request(t, "GET", "/profile?tab=liked", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, "A useful question") {
		t.Fatal("liked post missing from profile")
	}
	// Session records store a hash, not the bearer credential sent to the browser.
	var sessionToken string
	s.db.QueryRow("SELECT session_token FROM Sessions LIMIT 1").Scan(&sessionToken)
	var cookie *http.Cookie
	for _, c := range s.client.Jar.Cookies(s.base) {
		if c.Name == "session_id" {
			cookie = c
		}
	}
	if cookie == nil || sessionToken == cookie.Value {
		t.Fatal("missing or unhashed session")
	}
	r, b = s.form(t, "/logout", url.Values{})
	checkStatus(t, r, 303, b)
	s.db.QueryRow("SELECT COUNT(*) FROM Sessions").Scan(&count)
	if count != 0 {
		t.Fatal("logout did not invalidate server session")
	}
	r, b = s.request(t, "GET", "/profile", "", "", false)
	checkStatus(t, r, 303, b)
}

func TestValidationAndRequestSecurity(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	r, b := s.request(t, "POST", "/post", "title=hello", "application/x-www-form-urlencoded", false)
	checkStatus(t, r, 403, b)
	for _, values := range []url.Values{
		{"title": {"     "}, "content": {"valid content here"}, "category[]": {"1"}},
		{"title": {"Invalid topic"}, "content": {"valid content here"}, "category[]": {"999999"}},
		{"title": {"Too many topics"}, "content": {"valid content here"}, "category[]": {"1", "2", "3", "4"}},
	} {
		r, b = s.form(t, "/post", values)
		checkStatus(t, r, 400, b)
	}
	r, b = s.form(t, "/add_comment", url.Values{"post_id": {"999"}, "content": {"A reply"}})
	checkStatus(t, r, 404, b)
	r, b = s.form(t, "/add_comment", url.Values{"post_id": {"1"}, "content": {"  "}})
	checkStatus(t, r, 400, b)
	for _, payload := range []string{`{"postId":1,"action":"drop","type":"post"}`, `{"postId":1,"action":"like","type":"Users"}`, `{"postId":1,"action":"like","type":"post"} {}`} {
		r, b = s.request(t, "POST", "/like_dislike", payload, "application/json", true)
		checkStatus(t, r, 400, b)
	}
	r, b = s.request(t, "POST", "/like_dislike", `{"postId":999,"action":"like","type":"post"}`, "application/json", true)
	checkStatus(t, r, 404, b)
	r, b = s.request(t, "GET", "/missing", "", "", false)
	checkStatus(t, r, 404, b)
	r, b = s.request(t, "GET", "/healthz", "", "", false)
	checkStatus(t, r, 200, b)
	if r.Header.Get("Content-Security-Policy") == "" || r.Header.Get("X-Content-Type-Options") != "nosniff" {
		t.Fatal("missing security headers")
	}
	req, _ := http.NewRequest("POST", s.server.URL+"/logout", nil)
	req.Header.Set("X-CSRF-Token", s.csrf())
	req.Header.Set("Origin", "https://untrusted.example")
	r, err := s.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	r.Body.Close()
	if r.StatusCode != 403 {
		t.Fatal("cross-origin write allowed")
	}
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM Posts").Scan(&count)
	if count != 0 {
		t.Fatal("invalid submission left a partial post")
	}
}

func TestSessionsExpireAndSurviveRestart(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	// A new handler has no memory of login; the persistent session still resolves.
	req := httptest.NewRequest("GET", "/profile", nil)
	for _, c := range s.client.Jar.Cookies(s.base) {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	New(s.db, false).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("session lost on restart: %d", w.Code)
	}
	s.db.Exec("UPDATE Sessions SET expires_at=?", time.Now().UTC().Add(-time.Hour))
	r, b := s.request(t, "GET", "/profile", "", "", false)
	checkStatus(t, r, 303, b)
}

func TestConcurrentReactions(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	r, b := s.form(t, "/post", url.Values{"title": {"Concurrency test"}, "content": {"A useful concurrency test."}, "category[]": {"1"}})
	checkStatus(t, r, 303, b)
	var wg sync.WaitGroup
	errs := make(chan string, 12)
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("POST", s.server.URL+"/like_dislike", strings.NewReader(`{"postId":1,"action":"like","type":"post"}`))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-CSRF-Token", s.csrf())
			res, err := s.client.Do(req)
			if err != nil {
				errs <- err.Error()
				return
			}
			io.Copy(io.Discard, res.Body)
			res.Body.Close()
			if res.StatusCode != 200 {
				errs <- res.Status
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM Likes_Dislikes WHERE post_id=1").Scan(&count)
	if count != 0 {
		t.Fatalf("12 toggles must cancel, got %d reactions", count)
	}
}

func TestSeedAndSchemaAreIdempotent(t *testing.T) {
	s := newSite(t)
	if err := SeedDemo(s.db); err != nil {
		t.Fatal(err)
	}
	if err := SeedDemo(s.db); err != nil {
		t.Fatal(err)
	}
	if _, err := s.db.Exec(schema); err != nil {
		t.Fatal(err)
	}
	var posts, users int
	s.db.QueryRow("SELECT (SELECT COUNT(*) FROM Posts),(SELECT COUNT(*) FROM Users)").Scan(&posts, &users)
	if posts != 6 || users != 4 {
		t.Fatalf("seed duplicated data: %d posts, %d users", posts, users)
	}
	r, b := s.request(t, "GET", "/", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, "MayaDemo") {
		t.Fatal("seed not rendered")
	}
}

func TestPaginationAndMemberPrivacy(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	for i := 1; i <= 25; i++ {
		if _, err := s.db.Exec("INSERT INTO Posts(user_id,title,content) VALUES(1,?,?)", fmt.Sprintf("Discussion %02d", i), "A useful perspective to share."); err != nil {
			t.Fatal(err)
		}
	}
	r, b := s.request(t, "GET", "/", "", "", false)
	checkStatus(t, r, 200, b)
	if strings.Count(b, "class=\"post-card\"") != 20 || !strings.Contains(b, "Next conversations") {
		t.Fatal("first page was not bounded")
	}
	r, b = s.request(t, "GET", "/?page=2", "", "", false)
	checkStatus(t, r, 200, b)
	if strings.Count(b, "class=\"post-card\"") != 5 || !strings.Contains(b, "Previous") {
		t.Fatal("remaining posts missing")
	}
	r, b = s.request(t, "GET", "/post-details?post_id=1", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, "Discussion 01") {
		t.Fatal("old discussion could not be read directly")
	}
	s.db.Exec("INSERT INTO Users(username,email,password) VALUES('OtherUser','other@example.com','disabled')")
	s.db.Exec("INSERT INTO Likes_Dislikes(user_id,post_id,like_dislike) VALUES(2,1,1)")
	r, b = s.request(t, "GET", "/profile?user=2&tab=liked", "", "", false)
	checkStatus(t, r, 200, b)
	if strings.Contains(b, "Discussion 01") {
		t.Fatal("another member's liked posts were disclosed")
	}
}

func TestInvalidCookieRecoveryAndSecureFlags(t *testing.T) {
	s := newSite(t)
	req := httptest.NewRequest("GET", "/login", nil)
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "bad"})
	w := httptest.NewRecorder()
	New(s.db, true).ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatal("invalid cookie prevented rendering")
	}
	var token string
	for _, c := range w.Result().Cookies() {
		if c.Name == "csrf_token" {
			if !c.Secure || !c.HttpOnly || c.SameSite != http.SameSiteLaxMode {
				t.Fatal("cookie flags missing")
			}
			token = c.Value
		}
	}
	if len(token) != 64 || !strings.Contains(w.Body.String(), token) {
		t.Fatal("rendered CSRF token did not match replacement cookie")
	}
}
