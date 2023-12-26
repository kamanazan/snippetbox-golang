package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox.kamanazan.net/internal/models"
	"snippetbox.kamanazan.net/internal/validator"
)

// Define a snippetCreateForm struct to represent the form data and validation
// errors for the form fields. Note that all the struct fields are deliberately
// exported (i.e. start with a capital letter). This is because struct fields
// must be exported in order to be read by the html/template package when
// rendering the template.
type snippetCreateForm struct {
	Title   string `form:"title"`
	Content string `form:"content"`
	Expired int    `form:"expired"`
	// embed struct here, so snippetCreateForm "inherit" everything in Validator
	validator.Validator `form:-`
}

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
	data.Form = snippetCreateForm{
		Expired: 1,
	}
	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) createSnippetPost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := app.decodeFormData(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.StringNotEmpty(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.StringNotEmpty(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.StringInLimit(form.Title, 150), "title", "This field can not be more than 150 characters")
	validDuration := []int{1, 7, 365}
	form.CheckField(validator.ValueInRange(form.Expired, validDuration), "expired", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippet.Insert(form.Title, form.Content, form.Expired)
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
