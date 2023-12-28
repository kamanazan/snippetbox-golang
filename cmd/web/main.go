package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore" // New import
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq" // we alias this import to blank identifier because we only need its init() function so it is registered in database/sql

	"snippetbox.kamanazan.net/internal/models"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippet        *models.SnippetModel
	user           *models.UsersModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func openDB(dsn string) (*sql.DB, error) {
	/*
	   The sql.Open() function doesn’t actually create any connections, all it does is initialize the
	   pool for future use. Actual connections to the database are established lazily, as and when
	   needed for the first time. So to verify that everything is set up correctly we need to use the
	   db.Ping() method to create a connection and check for any errors.
	*/
	db, err_on_open := sql.Open("postgres", dsn)
	if err_on_open != nil {
		return nil, err_on_open
	}
	if err_timeout := db.Ping(); err_timeout != nil {
		return nil, err_timeout
	}

	return db, nil
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls. The value of the
	// flag will be stored in the addr variable at runtime.
	addr := flag.String("addr", ":4000", "Define adress:port")
	dsn := flag.String("dsn", "postgresql://kamanazan@localhost/snippet?sslmode=disable", "provide database connection string")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	// User can list available params from the '-help' param
	flag.Parse()

	//Define logging with its various prefix
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err_db_connection := openDB(*dsn)
	if err_db_connection != nil {
		errorLog.Fatal(err_db_connection)
	}

	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)

	templateCache, err_template := newTemplateCache()
	if err_template != nil {
		errorLog.Fatal(err_template)
	}

	formDecoder := form.NewDecoder()

	app := &application{ // the struct serve as dependeny injection, we defined it here and pass it to the handler function
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippet:        &models.SnippetModel{DB: db},
		user:           &models.UsersModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems. (p.71)
	srv := &http.Server{ // why use pointer ?
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// The value returned from the flag.String() function is a pointer to the flag
	// value, not the value itself. So we need to dereference the pointer (i.e.
	// prefix it with the * symbol) before using it. Note that we're using the
	// log.Printf() function to interpolate the address with the log message.
	infoLog.Printf("Starting server on %s", *addr)

	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
