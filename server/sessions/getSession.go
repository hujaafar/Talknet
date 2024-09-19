package sessions

import (
    "net/http"
)

func GetSessionUserID(r *http.Request) (int, bool) {
    cookie, err := r.Cookie("session_id")
    if err != nil {
        return -1, false
    }

    sessionID := cookie.Value
    userID, ok := sessionStore[sessionID]
    if !ok {
        return -1, false
    }

    // Verify that the session ID matches the one stored for the user
    if currentSessionID, exists := userSession[userID]; !exists || currentSessionID != sessionID {
        // The session is invalid or has been replaced
        return -1, false
    }

    return userID, true
}
