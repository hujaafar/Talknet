package sessions

import (
	"net/http"
	"time"
	"github.com/google/uuid"
)

var sessionStore = map[string]string{}

func CreateSession(w http.ResponseWriter, userID int) {
    sessionID := uuid.New().String()
    sessionStore[sessionID] = string(userID)

    http.SetCookie(w, &http.Cookie{
        Name:    "session_id",
        Value:   sessionID,
        Path:    "/",
        Expires: time.Now().Add(24 * time.Hour), // Example expiration
    })
}