package main

import (
	"BookBase/datalayer"
	"BookBase/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

func getBooksHandler(w http.ResponseWriter, r *http.Request) {
	limit, skip := extractPaging(r)

	var books = datalayer.GetBooks(limit, skip)
	renderJSON(w, books)
}

func extractPaging(r *http.Request) (*int, *int) {
	size := r.URL.Query().Get("size")
	page := r.URL.Query().Get("page")
	var limit *int
	var skip *int

	if len(size) != 0 && len(page) != 0 {
		sizeInt, err := strconv.Atoi(size)
		if err == nil {
			limit = new(int)
			*limit = sizeInt
		}
		pageInt, err := strconv.Atoi(page)
		if err == nil && limit != nil {
			s := (*limit) * (pageInt - 1)
			skip = new(int)
			*skip = s
		}
	}
	return limit, skip
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
	vars := mux.Vars(r)
	id := vars["id"]
	book, err := datalayer.GetBook(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	renderJSON(w, book)
}
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/book/{id}", getBookHandler)
	router.HandleFunc("/api/book/{id}/items", getBookItemsHandler)
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
