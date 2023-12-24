package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

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
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	// we avoid using http.HandleFunc directly because it comese from DefaultServeMux which is global.
	// so if we use package that access it and the package is compromised, attacker can inject malicious code.
	// so using localized HandleFunc like this is better for security
	// TODO: prove it!
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/create", app.create_snippet)
	mux.HandleFunc("/snippet/view", app.snippet_view)

	return mux
}
