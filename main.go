package main

import (
	"BookBase/datalayer"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/book/{id}", app.getBookHandler)
	router.HandleFunc("/api/book/{id}/items", app.getBookItemsHandler)
	router.HandleFunc("/api/book/{id}/giveOut", app.giveOutBookHandler).Methods("POST")
	router.HandleFunc("/api/returnBook", app.returnBookHandler).Methods("POST")
	router.HandleFunc("/api/book", app.getBooksHandler)
	router.HandleFunc("/api/author", app.getAuthorsHandler)
	router.HandleFunc("/api/author/{id:[0-9]+}", app.getAuthorHandler)

	http.Handle("/", router)

	err := datalayer.Init(app.infoLog)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	app.infoLog.Printf("Server is listening...")
	app.errorLog.Fatal(http.ListenAndServe("localhost:8181", nil))
}
