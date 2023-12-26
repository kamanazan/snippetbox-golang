package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox.kamanazan.net/internal/models"
)

// with  this all function here will be method for 'application' struct and have access
// to centralized logging.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippet.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData()
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	
    err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expired, err := strconv.Atoi(r.PostForm.Get("expired"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fieldErrors := map[string]string{}

	if strings.TrimSpace(title) == "" {
		fieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 150 {
		fieldErrors["title"] = "This field can not be more than 150 characters"
	}

	if strings.TrimSpace(content) == "" {
		fieldErrors["content"] = "This field can not be empty"
	}

	// Check the expires value matches one of the permitted values (1, 7 or
	// 365).
	if expired != 1 && expired != 7 && expired != 365 {
		fieldErrors["expired"] = "This field must equal 1, 7 or 365"
	}
	// If there are any errors, dump them in a plain text HTTP response and
	// return from the handler.
	if len(fieldErrors) > 0 {
		fmt.Fprint(w, fieldErrors)
		return
	}

	id, err := app.snippet.Insert(title, content, expired)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
    params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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
