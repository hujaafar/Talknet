package sessions

import (
	"net/http"
)

func GetSessionUserID(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return -1, false
	}
	userID, ok := sessionStore[cookie.Value]
	if !ok {
		return -1, false
	}

	return userID, true
}
