package forum

import (
	"net/url"
	"strings"
	"testing"
)

func TestErrorPagesRetainAuthenticatedNavigation(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	r, b := s.request(t, "GET", "/post-details?post_id=missing", "", "", false)
	checkStatus(t, r, 400, b)
	if !strings.Contains(b, "ReaderOne") || !strings.Contains(b, "Sign out") {
		t.Fatal("an error page lost the signed-in identity")
	}
	// This handler renders its error while its write transaction is still open.
	r, b = s.form(t, "/bookmarks", url.Values{"post_id": {"999"}, "action": {"save"}})
	checkStatus(t, r, 404, b)
	if !strings.Contains(b, "ReaderOne") {
		t.Fatal("transaction error lost the identity")
	}
}
