package forum

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuthCapacityStillAllowsExistingClientsWithinQuota(t *testing.T) {
	a := &App{attempts: authAttempts{entries: make(map[string]attempt)}}
	until := time.Now().Add(time.Minute)
	for i := 0; i < 10000; i++ {
		a.attempts.entries[fmt.Sprint(i)] = attempt{count: 1, until: until}
	}
	r := httptest.NewRequest("POST", "/login", nil)
	r.RemoteAddr = "0"
	if !a.allowAuth(r) {
		t.Fatal("an existing client lost its remaining quota")
	}
	r.RemoteAddr = "new-client"
	if a.allowAuth(r) {
		t.Fatal("a new client exceeded the memory bound")
	}
	r.RemoteAddr = "0"
	a.attempts.entries["0"] = attempt{count: 10, until: until}
	if a.allowAuth(r) {
		t.Fatal("an exhausted client exceeded its quota")
	}
}
