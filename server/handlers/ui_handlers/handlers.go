package ui_handlers

import "net/http"

func PageNotFound(w http.ResponseWriter, r *http.Request) {
	d := templateData{Title: "Page not found", DisableNavbar: true}
	renderTemplate(w, "404", d)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	d := templateData{Title: "Home"}
	renderTemplate(w, "home", d)
}
