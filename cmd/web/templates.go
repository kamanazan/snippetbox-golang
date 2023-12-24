package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.kamanazan.net/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
// TODO: why not use map ?
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		layout_files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			page,
		}
		// the following '...' is like destructuring in javascript
		ts, err := template.ParseFiles(layout_files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
