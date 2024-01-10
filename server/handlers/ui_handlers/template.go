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
		renderError(w)
	}
}

// renderError will render a generic internal server error page.
func renderError(w io.Writer) {
	// If we get to the point that we can't even render the error template, just do nothing
	_ = renderTemplateErr(w, "errors/500", templateData{
		Title:         "Error",
		DisableNavbar: true,
	})
}
