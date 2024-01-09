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
	validator.Validator `form:"-"`
}

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"pasword"`
	validator.Validator `form:"-"`
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

	err := app.decodeFormData(r, &form)
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

	app.sessionManager.Put(r.Context(), "flash", "Snippet Created")

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

	flash := app.sessionManager.PopString(r.Context(), "flash")

	data := app.newTemplateData()
	data.Snippet = snippet
	data.FlashMsg = flash

	app.render(w, http.StatusOK, "view.html", data)

}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = userSignUpForm{}

	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm

	err := app.decodeFormData(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.StringNotEmpty(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.StringNotEmpty(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.ValidEmail(form.Email), "email", "This field must be a valid email address")
	form.CheckField(validator.StringNotEmpty(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// If there are any errors, redisplay the signup form along with a 422
	// status code.
	if !form.Valid() {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}
	// Otherwise send the placeholder response (for now!).

	err = app.user.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		fmt.Printf("is error %s %v", err, errors.Is(err, models.ErrDuplicateEmail))
		if errors.Is(err, models.ErrDuplicateEmail) {
			data := app.newTemplateData()
			form.AddFieldError("email", "Email already exist")
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
			return

		}
	}
	http.Redirect(w, r, "/", http.StatusAccepted)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
