package forum

import (
	"net/url"
	"strings"
	"testing"
)

func TestShowcasePreservesExistingCommunityAndDoesNotDuplicate(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	createDiscussion(t, s, "An existing member discussion")
	r, b := s.form(t, "/add_comment", url.Values{"post_id": {"1"}, "content": {"An existing member reply."}})
	checkStatus(t, r, 303, b)
	var originalHash string
	if err := s.db.QueryRow("SELECT password FROM Users WHERE id=1").Scan(&originalHash); err != nil {
		t.Fatal(err)
	}
	summary, err := SeedShowcase(s.db)
	if err != nil {
		t.Fatal(err)
	}
	if summary.Users != 8 || summary.Posts != 24 || summary.Replies != 48 || summary.Likes != 84 || summary.AlreadyAdded {
		t.Fatalf("unexpected import: %+v", summary)
	}
	again, err := SeedShowcase(s.db)
	if err != nil || !again.AlreadyAdded {
		t.Fatalf("repeat import: %+v %v", again, err)
	}
	var users, posts, comments, topics int
	if err = s.db.QueryRow("SELECT (SELECT COUNT(*) FROM Users),(SELECT COUNT(*) FROM Posts),(SELECT COUNT(*) FROM Comments),(SELECT COUNT(DISTINCT category_id) FROM Post_Categories)").Scan(&users, &posts, &comments, &topics); err != nil {
		t.Fatal(err)
	}
	if users != 9 || posts != 25 || comments != 49 || topics != 12 {
		t.Fatalf("unexpected totals: %d %d %d %d", users, posts, comments, topics)
	}
	var title, hash string
	if err = s.db.QueryRow("SELECT p.title,u.password FROM Posts p JOIN Users u ON u.id=p.user_id WHERE p.id=1").Scan(&title, &hash); err != nil {
		t.Fatal(err)
	}
	if title != "An existing member discussion" || hash != originalHash {
		t.Fatal("existing records changed")
	}
	var malformed int
	if err = s.db.QueryRow("SELECT COUNT(*) FROM Users WHERE id<>1 AND (username NOT LIKE '%Demo' OR email NOT LIKE '%@example.invalid')").Scan(&malformed); err != nil {
		t.Fatal(err)
	}
	if malformed != 0 {
		t.Fatal("sample identities are not clearly fictional")
	}
	r, b = s.request(t, "GET", "/?page=2", "", "", false)
	checkStatus(t, r, 200, b)
	if strings.Count(b, `class="post-card"`) != 5 {
		t.Fatal("showcase did not populate a second feed page")
	}
}

func TestShowcaseCollisionRollsBackInsteadOfAdoptingAnAccount(t *testing.T) {
	s := newSite(t)
	// The second author collides, so the first insertion must also roll back.
	if _, err := s.db.Exec("INSERT INTO Users(username,email,password) VALUES('TheoDemo','existing@example.invalid','existing-password-hash')"); err != nil {
		t.Fatal(err)
	}
	if _, err := SeedShowcase(s.db); err == nil {
		t.Fatal("showcase adopted an existing account")
	}
	var users, posts, markers int
	if err := s.db.QueryRow("SELECT (SELECT COUNT(*) FROM Users),(SELECT COUNT(*) FROM Posts),(SELECT COUNT(*) FROM Demo_Seeds)").Scan(&users, &posts, &markers); err != nil {
		t.Fatal(err)
	}
	if users != 1 || posts != 0 || markers != 0 {
		t.Fatal("a failed import left partial data")
	}
}
