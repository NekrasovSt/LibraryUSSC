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
func GetBook(isbn string) (models.Book, error) {
	var book models.Book
	result := db.First(&book, "isbn = ?", isbn)
	if result.Error != nil {
		return models.Book{}, result.Error
	}
	return book, nil
}
