package forum

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func signInSecondAuthor(t *testing.T, s *testSite) {
	t.Helper()
	r, b := s.form(t, "/logout", url.Values{})
	checkStatus(t, r, 303, b)
	r, b = s.form(t, "/register", url.Values{"username": {"ReaderTwo"}, "email": {"second@example.invalid"}, "password": {"Another-long-passphrase-456"}})
	checkStatus(t, r, 303, b)
	r, b = s.form(t, "/login", url.Values{"username": {"ReaderTwo"}, "password": {"Another-long-passphrase-456"}})
	checkStatus(t, r, 303, b)
}
func createDiscussion(t *testing.T, s *testSite, title string) {
	t.Helper()
	r, b := s.form(t, "/post", url.Values{"title": {title}, "content": {"A considered perspective for the community to explore."}, "category[]": {"1"}})
	checkStatus(t, r, 303, b)
}
func bookmarkJSON(t *testing.T, s *testSite, id, action string) (*http.Response, string) {
	t.Helper()
	req, err := http.NewRequest("POST", s.server.URL+"/bookmarks", strings.NewReader(url.Values{"post_id": {id}, "action": {action}}.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CSRF-Token", s.csrf())
	r, err := s.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	return r, string(b)
}
func TestBookmarksArePrivatePersistentAndIdempotent(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	createDiscussion(t, s, "An unusual Alpha perspective")
	signInSecondAuthor(t, s)
	createDiscussion(t, s, "The author's own Bravo discussion")
	for range 2 {
		r, b := bookmarkJSON(t, s, "1", "save")
		checkStatus(t, r, 200, b)
		if !strings.Contains(b, `"saved":true`) {
			t.Fatal(b)
		}
	}
	var count int
	s.db.QueryRow("SELECT COUNT(*) FROM Bookmarks WHERE user_id=2 AND post_id=1").Scan(&count)
	if count != 1 {
		t.Fatalf("repeated save produced %d rows", count)
	}
	r, b := s.request(t, "GET", "/profile?tab=saved", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, "An unusual Alpha perspective") || strings.Contains(b, "The author&#39;s own Bravo discussion") {
		t.Fatal("saved feed did not load the saved post alone")
	}
	r, b = s.form(t, "/logout", url.Values{})
	checkStatus(t, r, 303, b)
	r, b = s.form(t, "/login", url.Values{"username": {"ReaderOne"}, "password": {"A-long-passphrase-123"}})
	checkStatus(t, r, 303, b)
	r, b = s.request(t, "GET", "/profile?user=2&tab=saved", "", "", false)
	checkStatus(t, r, 200, b)
	if strings.Contains(b, "An unusual Alpha perspective") || !strings.Contains(b, "Bravo discussion") {
		t.Fatal("another member could view private bookmarks")
	}
	r, b = bookmarkJSON(t, s, "999", "save")
	checkStatus(t, r, 404, b)
	r, b = bookmarkJSON(t, s, "1", "toggle")
	checkStatus(t, r, 400, b)
	r, b = s.form(t, "/logout", url.Values{})
	checkStatus(t, r, 303, b)
	r, b = s.form(t, "/login", url.Values{"username": {"ReaderTwo"}, "password": {"Another-long-passphrase-456"}})
	checkStatus(t, r, 303, b)
	for range 2 {
		r, b = bookmarkJSON(t, s, "1", "remove")
		checkStatus(t, r, 200, b)
	}
	s.db.QueryRow("SELECT COUNT(*) FROM Bookmarks").Scan(&count)
	if count != 0 {
		t.Fatal("bookmark removal failed")
	}
	r, b = s.form(t, "/logout", url.Values{})
	checkStatus(t, r, 303, b)
	r, b = bookmarkJSON(t, s, "1", "save")
	checkStatus(t, r, 401, b)
}

func TestEditingRequiresOwnershipAndRejectsStaleVersions(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	createDiscussion(t, s, "The original discussion")
	r, b := s.form(t, "/add_comment", url.Values{"post_id": {"1"}, "content": {"An existing reply worth preserving."}})
	checkStatus(t, r, 303, b)
	r, b = s.request(t, "GET", "/post/edit?post_id=1", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, "The original discussion") || !strings.Contains(b, `name="revision"`) {
		t.Fatal("editor did not load the existing discussion")
	}
	values := url.Values{"post_id": {"1"}, "revision": {"1"}, "title": {"The revised discussion"}, "content": {"A clearer perspective with enough context for readers."}, "category[]": {"2"}}
	r, b = s.form(t, "/post/edit", values)
	checkStatus(t, r, 303, b)
	var title string
	var comments, category, revision int
	s.db.QueryRow("SELECT title FROM Posts WHERE id=1").Scan(&title)
	s.db.QueryRow("SELECT COUNT(*) FROM Comments WHERE post_id=1").Scan(&comments)
	s.db.QueryRow("SELECT category_id FROM Post_Categories WHERE post_id=1").Scan(&category)
	s.db.QueryRow("SELECT revision FROM Post_Revisions WHERE post_id=1").Scan(&revision)
	if title != "The revised discussion" || comments != 1 || category != 2 || revision != 2 {
		t.Fatal("editing failed to update atomically while preserving replies")
	}
	values.Set("title", "A conflicting tab's changes")
	r, b = s.form(t, "/post/edit", values)
	checkStatus(t, r, 409, b)
	if !strings.Contains(b, "A conflicting tab&#39;s changes") || !strings.Contains(b, "changed in another tab") {
		t.Fatal("conflict did not retain the submitted text")
	}
	s.db.QueryRow("SELECT title FROM Posts WHERE id=1").Scan(&title)
	if title != "The revised discussion" {
		t.Fatal("stale edit overwrote the current version")
	}
	signInSecondAuthor(t, s)
	r, b = s.request(t, "GET", "/post/edit?post_id=1", "", "", false)
	checkStatus(t, r, 403, b)
	values.Set("revision", "2")
	r, b = s.form(t, "/post/edit", values)
	checkStatus(t, r, 403, b)
	s.db.QueryRow("SELECT title FROM Posts WHERE id=1").Scan(&title)
	if title != "The revised discussion" {
		t.Fatal("another author edited this discussion")
	}
}

func TestQuickSearchReturnsOnlyPublicMetadata(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	createDiscussion(t, s, "An astronomy question")
	r, b := s.request(t, "GET", "/api/search?q=astronomy", "", "", false)
	checkStatus(t, r, 200, b)
	var payload struct {
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal([]byte(b), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Results) != 1 || len(payload.Results[0]) != 4 || payload.Results[0]["title"] != "An astronomy question" {
		t.Fatal("unexpected search metadata")
	}
	if strings.Contains(b, "reader@example.com") || strings.Contains(b, "password") {
		t.Fatal("search leaked private account data")
	}
	r, b = s.request(t, "GET", "/api/search?q=a", "", "", false)
	checkStatus(t, r, 200, b)
	if !strings.Contains(b, `"results":[]`) {
		t.Fatal("short search should be empty")
	}
	r, b = s.request(t, "GET", "/api/search?q="+strings.Repeat("a", 201), "", "", false)
	checkStatus(t, r, 400, b)
}
