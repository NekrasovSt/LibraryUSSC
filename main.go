package main

import (
	"BookBase/datalayer"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/book/{id}", getBookHandler)
	router.HandleFunc("/api/book/{id}/items", getBookItemsHandler)
	router.HandleFunc("/api/book", getBooksHandler)
	router.HandleFunc("/api/author", getAuthorsHandler)
	router.HandleFunc("/api/author/{id:[0-9]+}", getAuthorHandler)

	http.Handle("/", router)

	err := datalayer.Init()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println("Server is listening...")
	http.ListenAndServe("localhost:8181", nil)
}
