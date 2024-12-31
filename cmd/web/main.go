package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/joho/godotenv"
	"github.com/uuzay/snippetbox/pkg/models/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		errorLog.Fatal("Error loading secrets")
	}

	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsnDefault := fmt.Sprintf("%v:%v@/%v?parseTime=true", dbUsername, dbPassword, dbName)

	secretDefault := os.Getenv("SECRET_KEY")

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", dsnDefault, "MySQL data source name")
	secret := flag.String("secret", secretDefault, "Secret key")
	flag.Parse()

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Starting server on :4000")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

// Prevent listing directories by sending requests to "/static/" etc.
func neuter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
