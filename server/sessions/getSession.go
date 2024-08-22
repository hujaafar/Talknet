package sessions

import (
	"net/http"
)

func GetSessionUserID(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", false
	}

	userID, ok := sessionStore[cookie.Value]
	if !ok {
		return "", false
	}

	return userID, true
}
