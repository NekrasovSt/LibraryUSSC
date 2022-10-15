package main

import (
	"BookBase/datalayer"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (app *application) getBooksHandler(w http.ResponseWriter, r *http.Request) {
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
func (app *application) getBookItemsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	items, err := datalayer.GetBookItems(id)
	if err != nil {
		app.handleError(err, w)
		return
	}
	renderJSON(w, items)
}
func renderJSON(w http.ResponseWriter, data interface{}) {
	js, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (app *application) handleError(err error, w http.ResponseWriter) {
	app.errorLog.Print(err)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (app *application) getBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	book, err := datalayer.GetBook(id)
	if err != nil {
		app.handleError(err, w)
		return
	}
	renderJSON(w, book)
}
func (app *application) getAuthorsHandler(w http.ResponseWriter, r *http.Request) {
	limit, skip := extractPaging(r)

	books, err := datalayer.GetAuthors(limit, skip)
	if err != nil {
		app.handleError(err, w)
		return
	}
	renderJSON(w, books)
}
func (app *application) getAuthorHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idParam := vars["id"]
	id, _ := strconv.Atoi(idParam)
	book, err := datalayer.GetAuthor(id)
	if err != nil {
		app.handleError(err, w)
		return
	}
	renderJSON(w, book)
}

type ItemDto struct {
	Id int `json:"id"`
}

func (app *application) giveOutBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	isbn := vars["id"]
	_, err := datalayer.GetBook(isbn)
	if err != nil {
		app.handleError(err, w)
		return
	}
	account := new(ItemDto)
	err = json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var id int
	id, err = datalayer.GiveOutBook(isbn, account.Id)
	if err != nil {
		app.handleError(err, w)
		return
	}
	renderJSON(w, id)
}
func (app *application) returnBookHandler(w http.ResponseWriter, r *http.Request) {
	bookItem := new(ItemDto)
	err := json.NewDecoder(r.Body).Decode(&bookItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = datalayer.ReturnBook(bookItem.Id)
	if err != nil {
		app.handleError(err, w)
		return
	}
}
