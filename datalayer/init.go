package datalayer

import (
	"BookBase/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}
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
	if newDB {
		err = fillExamples(conn)
		if err != nil {
			return err
		}
	}
	return nil
}

func fillExamples(conn *gorm.DB) error {
	authors := models.GetDefaultAuthors()
	err := conn.Create(&authors).Error
	if err != nil {
		return err
	}

	publishers := models.GetDefaultPublishers()
	err = conn.Create(&publishers).Error
	if err != nil {
		return err
	}

	users := models.GetDefaultUsers()
	err = conn.Create(&users).Error
	if err != nil {
		return err
	}

	books := models.GetDefaultBooks()

	books[0].Authors = append(books[0].Authors, authors[0], authors[1])
	books[0].PublisherId = publishers[0].Id

	books[0].Books = models.GetBookItems()
	err = conn.Create(&books).Error
	if err != nil {
		return err
	}
	return nil
}
