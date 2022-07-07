package main

import (
	"BookBase/datalayer"
	"BookBase/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, models.GetDefaultBooks())
}
func getBookItemsHandler(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, models.GetBookItems())
}
func renderJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func getBookHandler(w http.ResponseWriter, r *http.Request) {
	renderJSON(w, models.GetDefaultBooks()[0])
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/book/{id:[0-9]+}", getBookHandler)
	router.HandleFunc("/api/book/{id:[0-9]+}/items", getBookItemsHandler)
	router.HandleFunc("/api/book", getBooksHandler)
	http.Handle("/", router)

	err := datalayer.Init()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println("Server is listening...")
	http.ListenAndServe("localhost:8181", nil)
}
