package utils

import (
	"encoding/json"
	"html/template"
	"net/http"
)

type GeneratedPage struct {
	Title string
	Src   string
}

func Json200(payload any, w http.ResponseWriter) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func ValidateAdmin(token string) bool {
	// Simple validation for now
	return token == "Bearer - admin"
}

func WriteImageHTML(title string, src string, w http.ResponseWriter) error {
	tmpl, err := template.New("imagePage").Parse(ImagePageHtml)
	if err != nil {
		return err
	}

	page := GeneratedPage{Title: title, Src: src}

	return tmpl.Execute(w, page)
}
