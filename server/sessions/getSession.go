package sessions

import (
	"net/http"
	"strconv"
)

func GetSessionUserID(r *http.Request) (int, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0, false
	}

	userID, ok := sessionStore[cookie.Value]
	if !ok {
		return 0, false
	}
	id, err := strconv.Atoi(userID)
	if err != nil {
		return 0, false
	}
	return id, true
}
