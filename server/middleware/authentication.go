package middleware

import (
	"net/http"
	"talknet/server/sessions"
)
func authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, loggedIn := sessions.GetSessionUserID(r)
        if !loggedIn {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}
