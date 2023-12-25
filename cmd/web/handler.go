package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"snippetbox.kamanazan.net/internal/models"
)

// with  this all function here will be method for 'application' struct and have access
// to centralized logging.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		// alternative of w.WriteHeader() and w.Write()
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) create_snippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		/*
		   By default Write will return 200, so if we want to send code other than 200 we should
		   call WriteHeader(http_code). WriteHeader can only be called once and should be done before Write().
		   Calling WriteHeader again will give warning <http: superfluous response.WriteHeader call from main.create_snippet>
		*/
		// w.WriteHeader(405)
		// w.Write([]byte("Method not Allowed"))

		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 365

	id, err := app.snippet.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) snippet_view(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippet.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData()
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.html", data)

}
