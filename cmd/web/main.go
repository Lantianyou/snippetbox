package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"lantianyou.com/snippetbox/pkg/models/mysql"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main(){

	addr := flag.String("addr", ":4000", "http address")
	dsn  := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name and password")
	flag.Parse()

	// Must be called after all flags are defined
	// and before flags are accessed by the program.

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate | log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate | log.Ltime | log.Lshortfile)

	db, err := OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplatesCache("./ui/html/")
	if err != nil {
		errorLog.Fatalln(err)
	}

	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &mysql.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("listening on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func OpenDB(dsn string) (*sql.DB, error){
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping();err!=nil {
		return nil, err
	}
	return db, err
}