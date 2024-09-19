package sessions

import (
    "net/http"
    "time"

    "github.com/google/uuid"
)

var sessionStore = map[string]int{}  // Maps sessionID to userID
var userSession = map[int]string{}   // Maps userID to sessionID

func CreateSession(w http.ResponseWriter, userID int) {
    // Check if the user already has an active session
    if oldSessionID, exists := userSession[userID]; exists {
        // Invalidate the old session
        delete(sessionStore, oldSessionID)
    }

    // Create a new session
    sessionID := uuid.New().String()
    sessionStore[sessionID] = userID
    userSession[userID] = sessionID

    // Set the session cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "session_id",
        Value:    sessionID,
        Path:     "/",
        Expires:  time.Now().Add(24 * time.Hour),
        HttpOnly: true, // Enhances security by preventing JavaScript access
        Secure:   true, // Ensures the cookie is sent over HTTPS
    })
}
