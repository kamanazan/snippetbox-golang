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

		/*
			old way to register all templates where it only assume we only have nav.html in html/partials folder.
			it will save effort if we can automatically add any new file inside html/partial.
			what we are going to do:
			1. parse base.html file
			2. since there *might* be more files in html/partial we are going to use parseGlob instead
			3. we already knew which file to add in html/pages so we just use parse file

		*/
		// layout_files := []string{
		// 	"./ui/html/base.html",
		// 	"./ui/html/partials/nav.html",
		// 	page,
		// }
		// // the following '...' is like destructuring in javascript
		// ts, err := template.ParseFiles(layout_files...)
		// if err != nil {
		// 	return nil, err
		// }

		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		// mistake use "template" instead of "ts" causing other templates not rendered
		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
