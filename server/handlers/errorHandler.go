package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"talknet/structs"
)

func RenderErrorPage(w http.ResponseWriter, errMsg string, statusCode int) {
	data := structs.ErrorData{
		ErrorMessage: errMsg,
		Code:         strconv.Itoa(statusCode),
	}

	tmpl, err := template.ParseFiles("static/pages/error.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	err = tmpl.Execute(w, data)
	if err != nil {
		// If template execution fails, log the error and respond with a simple error message
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
