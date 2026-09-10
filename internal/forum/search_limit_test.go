package forum

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCommandSearchReturnsOnlySixNewestMatches(t *testing.T) {
	s := newSite(t)
	s.signup(t)
	for i := 1; i <= 8; i++ {
		createDiscussion(t, s, fmt.Sprintf("Astronomy discussion %d", i))
	}
	r, b := s.request(t, "GET", "/api/search?q=astronomy", "", "", false)
	checkStatus(t, r, 200, b)
	var payload struct {
		Results []struct {
			ID int `json:"id"`
		} `json:"results"`
	}
	if err := json.Unmarshal([]byte(b), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Results) != 6 || payload.Results[0].ID != 8 || payload.Results[5].ID != 3 {
		t.Fatal("search did not return the six newest matches")
	}
}
