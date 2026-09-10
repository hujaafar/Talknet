package forum

import "testing"

func TestProfileTargetsArePositiveAndExplicitUserTakesPrecedence(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	for _, query := range []string{"user=0", "user=-1", "user=broken", "id=0", "id=-1", "id=broken"} {
		r, b := s.request(t, "GET", "/profile?"+query, "", "", false)
		checkStatus(t, r, 400, b)
	}
	r, b := s.request(t, "GET", "/profile?user=1&id=999", "", "", false)
	checkStatus(t, r, 200, b)
}
