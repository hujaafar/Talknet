package middleware

import (
	"net/http"
	"talknet/server/handlers"
	"talknet/server/sessions"
)
func authenticate(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        _, loggedIn := sessions.GetSessionUserID(r)
        if !loggedIn {
           handlers.RenderErrorPage(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}
