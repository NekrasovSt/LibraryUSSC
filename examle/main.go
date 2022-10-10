package main

import (
	"BookBase/models"
	"encoding/csv"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *gorm.DB
var publishers map[string]models.Publisher
var authors map[string]models.Author

func Init(log *log.Logger) error {
	newLogger := logger.New(
		log,
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		return e
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, "postgres", password)
	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{Logger: newLogger})
	if err != nil {
		return err
	}
	var bases []string
	conn.Raw(fmt.Sprintf("SELECT \"datname\" FROM pg_database WHERE datname = '%s'", dbName)).Scan(&bases)
	newDB := len(bases) == 0
	if newDB {
		_ = conn.Exec("CREATE DATABASE " + dbName)
	}

	dbUri = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	conn, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true, Logger: newLogger})
	if err != nil {
		return err
	}
	modelsForMigration := []interface{}{
		&models.Book{},
		&models.BookItem{},
		&models.Publisher{},
		&models.User{},
		&models.Author{},
	}
	db = conn
	err = db.AutoMigrate(modelsForMigration...)
	if err != nil {
		return err
	}
	conn, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: false, Logger: newLogger})
	if err != nil {
		return err
	}
	db = conn
	err = db.AutoMigrate(modelsForMigration...)
	if err != nil {
		return err
	}

	if !newDB {
		return nil
	}
	err = fillExamples(log, conn)
	if err != nil {
		return err
	}
	return nil
}
func getOrCreatePublisher(conn *gorm.DB, name string) (models.Publisher, error) {
	if val, ok := publishers[name]; ok {
		return val, nil
	}
	var book models.Publisher
	result := conn.First(&book, "name = ?", name)
	if result.Error == gorm.ErrRecordNotFound {
		newPublisher := models.Publisher{
			Name:  name,
			Email: "some@mmail.com",
		}
		err := conn.Create(&newPublisher).Error
		if err != nil {
			return models.Publisher{}, err
		}
		publishers[name] = newPublisher
		return newPublisher, nil
	}
	if result.Error != nil {
		return models.Publisher{}, result.Error
	}
	publishers[name] = book
	return book, nil
}
func getOrCreateAuthor(conn *gorm.DB, name string) (models.Author, error) {
	if val, ok := authors[name]; ok {
		return val, nil
	}
	pair := strings.Split(name, " ")
	var author models.Author
	firstName := pair[0]
	var secondName string
	if len(pair) == 1 {
		secondName = "Anonimus"
	} else {
		secondName = pair[1]
	}
	result := conn.First(&author, "first_name = ? and second_name = ?", firstName, secondName)
	if result.Error == gorm.ErrRecordNotFound {
		newAuthor := models.Author{
			FirstName:  firstName,
			SecondName: secondName,
		}
		err := conn.Create(&newAuthor).Error
		if err != nil {
			return models.Author{}, err
		}
		authors[name] = newAuthor
		return newAuthor, nil
	}
	if result.Error != nil {
		return models.Author{}, result.Error
	}
	authors[name] = author
	return author, nil
}
func fillExamples(errorLog *log.Logger, conn *gorm.DB) error {
	args := os.Args
	if len(args) < 2 {
		errorLog.Fatal("Нужно указать название файла")
	}
	filePath := args[1]
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	csvReader.Comma = ','
	csvReader.LazyQuotes = true
	// skip first
	csvReader.Read()
	return conn.Transaction(func(tx *gorm.DB) error {
		for {
			record, e := csvReader.Read()
			if e == io.EOF {
				break
			}
			if e != nil {
				errorLog.Print(e)
				continue
			}
			publisherName := record[11]
			items := strings.Split(record[10], "/")
			y, err := strconv.Atoi(items[2])
			if err != nil {
				return err
			}
			m, err := strconv.Atoi(items[0])
			if err != nil {
				return err
			}
			d, err := strconv.Atoi(items[1])
			if err != nil {
				return err
			}
			publicationDate := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
			isbn := record[5]
			author := record[2]

			title := record[1]
			pub, err := getOrCreatePublisher(tx, publisherName)
			if err != nil {
				return err
			}
			aut, err := getOrCreateAuthor(tx, author)
			if err != nil {
				return err
			}
			newBook := models.Book{
				PublisherId: pub.Id,
				Year:        publicationDate.Year(),
				ISBN:        isbn,
				Title:       title,
			}
			newBook.Authors = append(newBook.Authors, aut)
			err = tx.Create(&newBook).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func main() {
	// https://www.kaggle.com/datasets/jealousleopard/goodreadsbooks?resource=download
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	publishers = make(map[string]models.Publisher)
	authors = make(map[string]models.Author)
	err := Init(infoLog)
	if err != nil {
		errorLog.Fatal(err.Error())
	}
}
