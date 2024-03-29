package datalayer

import (
	"BookBase/models"
	"time"
)

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
func GetUsers(limit, skip *int) ([]models.User, error) {
	var users []models.User
	conn := GetDB()
	var err error
	if limit != nil && skip != nil {
		err = conn.Limit(*limit).Offset(*skip).Find(&users).Error
	} else {
		err = conn.Find(&users).Error
	}
	return users, err
}
func GetAuthors(limit, skip *int) ([]models.Author, error) {
	var authors []models.Author
	conn := GetDB()
	var err error
	if limit != nil && skip != nil {
		err = conn.Limit(*limit).Offset(*skip).Find(&authors).Error
	} else {
		err = conn.Find(&authors).Error
	}
	return authors, err
}
func GetBook(isbn string) (models.Book, error) {
	var book models.Book
	result := db.First(&book, "isbn = ?", isbn)
	if result.Error != nil {
		return models.Book{}, result.Error
	}
	return book, nil
}
func ReturnBook(bookId int) error {
	var bookItem models.BookItem
	err := db.Where("id = ? and user_id is not null", bookId).First(&bookItem).Error
	if err != nil {
		return err
	}
	bookItem.UserId = nil
	return db.Save(bookItem).Error
}
func AddBookItems(isbn string, count int) ([]models.BookItem, error) {
	var book models.Book
	result := db.First(&book, "isbn = ?", isbn)
	if result.Error != nil {
		return nil, result.Error
	}
	bookItems := []models.BookItem{}
	for i := 0; i < count; i++ {
		bookItems = append(bookItems, models.BookItem{ISBN: isbn, Receipt: time.Now().UTC()})
	}
	err := db.Create(&bookItems).Error
	if err != nil {
		return nil, result.Error
	}
	return bookItems, nil
}
func GiveOutBook(isbn string, userId int) (int, error) {
	var bookItem models.BookItem
	err := db.Where("isbn = ? and user_id is null", isbn).First(&bookItem).Error
	if err != nil {
		return 0, err
	}
	bookItem.UserId = &userId
	return bookItem.Id, db.Save(&bookItem).Error
}
func CreateUser(user *models.User) error {
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetAuthor(id int) (models.Author, error) {
	var author models.Author
	result := db.First(&author, id)
	if result.Error != nil {
		return models.Author{}, result.Error
	}
	return author, nil
}
func GetBookItems(isbn string) ([]models.BookItem, error) {
	var bookItems []models.BookItem
	err := db.Where("isbn = ?", isbn).Find(&bookItems).Error
	return bookItems, err
}
