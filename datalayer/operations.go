package datalayer

import "BookBase/models"

func GetBooks(limit, skip *int) []models.Book {
	var books []models.Book
	conn := GetDB()
	if limit != nil && skip != nil {
		conn.Limit(*limit).Offset(*skip).Find(&books)
	} else {
		conn.Find(&books)
	}
	return books
}
func GetAuthors(limit, skip *int) []models.Author {
	var authors []models.Author
	conn := GetDB()
	if limit != nil && skip != nil {
		conn.Limit(*limit).Offset(*skip).Find(&authors)
	} else {
		conn.Find(&authors)
	}
	return authors
}
func GetBook(isbn string) (models.Book, error) {
	var book models.Book
	result := db.First(&book, "isbn = ?", isbn)
	if result.Error != nil {
		return models.Book{}, result.Error
	}
	return book, nil
}
func GetAuthor(id int) (models.Author, error) {
	var author models.Author
	result := db.First(&author, id)
	if result.Error != nil {
		return models.Author{}, result.Error
	}
	return author, nil
}
func GetBookItems(isbn string) []models.BookItem {
	var bookItems []models.BookItem
	db.Where("isbn = ?", isbn).Find(&bookItems)
	return bookItems
}
