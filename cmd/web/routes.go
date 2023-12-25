package main

import (
    "net/http"
    "github.com/justinas/alice"
    "github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
    router := httprouter.New()
    // Create a handler function which wraps our notFound() helper, and then
    // assign it as the custom handler for 404 Not Found responses. You can also
    // set a custom handler for 405 Method Not Allowed responses by setting
    // router.MethodNotAllowed in the same way too.
    router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.notFound(w)
    })
	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	/*
	   When this handler receives a request, it will remove the leading slash from the URL path and
	   then search the ./ui/static directory for the corresponding file to send to the user.
	   So, for this to work correctly, we must strip the leading "/static" from the URL path before
	   passing it to http.FileServer.
	   Otherwise it will be looking for a file which doesn’t exist and
	   the user will receive a 404 page not found response.
	   Fortunately Go includes a http.StripPrefix() helper specifically for this task.
	*/
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.createSnippet)
    router.HandlerFunc(http.MethodPost, "/snippet/create", app.createSnippetPost)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.viewSnippet)

    middlewares := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return  middlewares.Then(router)
}
