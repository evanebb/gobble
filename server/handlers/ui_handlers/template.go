package ui_handlers

import (
	"github.com/evanebb/gobble/resources"
	"html/template"
	"io"
)

type templateData struct {
	Title         string
	Data          any
	DisableNavbar bool
}

// renderTemplateErr will render the template referenced by the passed name and return an error if it encounters one
func renderTemplateErr(w io.Writer, templateName string, d templateData) error {
	var err error

	templateFile := "templates/" + templateName + ".gohtml"
	tmpl, err := template.New("base").ParseFS(resources.Templates, "templates/base.gohtml", templateFile)
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, d)
	if err != nil {
		return err
	}

	return nil
}

// renderTemplate will render the template referenced by the passed name and render an error page if an error occurs
func renderTemplate(w io.Writer, templateName string, d templateData) {
	err := renderTemplateErr(w, templateName, d)
	if err != nil {
		renderError(w, err)
	}
}

// renderError will render the error template with the message from the given error
func renderError(w io.Writer, err error) {
	data := templateData{
		Title: "Error",
		Data:  err.Error(),
	}

	// If we get to the point that we can't even render the error template, just do nothing
	_ = renderTemplateErr(w, "error", data)
}
