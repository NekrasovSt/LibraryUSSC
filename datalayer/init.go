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

type AppConfig struct {
	DbName   string
	DbPort   string
	Password string
	UserName string
	DbHost   string
	Endpoint string
}

func ParseConfig() (*AppConfig, error) {
	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		return nil, e
	}
	config := new(AppConfig)
	config.UserName = os.Getenv("db_user")
	config.Password = os.Getenv("db_pass")
	config.DbName = os.Getenv("db_name")
	config.DbHost = os.Getenv("db_host")
	config.DbPort = os.Getenv("db_port")
	config.Endpoint = os.Getenv("endpoint")
	return config, nil
}

func Init(config *AppConfig, log *log.Logger) error {
	newLogger := logger.New(
		log,
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", config.DbHost, config.UserName, "postgres", config.Password, config.DbPort)
	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{Logger: newLogger})
	if err != nil {
		return err
	}
	var bases []string
	conn.Raw(fmt.Sprintf("SELECT \"datname\" FROM pg_database WHERE datname = '%s'", config.DbName)).Scan(&bases)
	newDB := len(bases) == 0
	if newDB {
		_ = conn.Exec("CREATE DATABASE " + config.DbName)
	}

	dbUri = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", config.DbHost, config.UserName, config.DbName, config.Password, config.DbPort)
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
	dialector := &postgres.Dialector{Config: &postgres.Config{DSN: dbUri, PreferSimpleProtocol: true}}
	conn, err = gorm.Open(dialector, &gorm.Config{DisableForeignKeyConstraintWhenMigrating: false, Logger: newLogger})
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
