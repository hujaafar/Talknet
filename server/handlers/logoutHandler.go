package handlers

import (
	"net/http"
	"time"
)


func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:    "session_id",
        Value:   "",
        Path:    "/",
        Expires: time.Now().Add(-time.Hour),
    })
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
