package sessions

import (
    "net/http"
)

func LogoutUser(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("session_id")
    if err == nil {
        sessionID := cookie.Value
        if userID, ok := sessionStore[sessionID]; ok {
            // Delete the session mappings
            delete(sessionStore, sessionID)
            delete(userSession, userID)
        }
        // Remove the session cookie
        http.SetCookie(w, &http.Cookie{
            Name:     "session_id",
            Value:    "",
            Path:     "/",
            MaxAge:   -1,     // Deletes the cookie immediately
            HttpOnly: true,
            Secure:   true,
        })
    }
}
