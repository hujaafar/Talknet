package forum

import (
	"net/url"
	"strings"
	"testing"
)

func TestSearchLimitCountsUnicodeCharacters(t *testing.T) {
	s := newSite(t)
	for _, route := range []string{"/?q=", "/api/search?q="} {
		for _, letters := range []int{200, 201} {
			r, b := s.request(t, "GET", route+url.QueryEscape(strings.Repeat("界", letters)), "", "", false)
			want := 200
			if letters > 200 {
				want = 400
			}
			checkStatus(t, r, want, b)
		}
	}
}
