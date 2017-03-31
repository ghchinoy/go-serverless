package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// parseTemplate applies a given file to the body of the base template.
func parseTemplate(filename string) *htmlTemplate {
	tmpl := template.Must(template.ParseFiles("templates/base.html"))

	// Put the named file into a template called "body"
	path := filepath.Join("templates", filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not read template: %v", err))
	}
	template.Must(tmpl.New("body").Parse(string(b)))

	return &htmlTemplate{tmpl.Lookup("base.html")}
}

type htmlTemplate struct {
	t *template.Template
}

func (tmpl *htmlTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}) *apiError {
	d := struct {
		Data interface{}
	}{
		Data: data,
	}

	if err := tmpl.t.Execute(w, d); err != nil {
		return apiErrorf(err, "couldn't write template: %v")
	}
	return nil
}
