package datalayer

import (
	"BookBase/models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}
func Init() error {
	e := godotenv.Load() //Загрузить файл .env
	if e != nil {
		return e
	}
	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, "postgres", password)
	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		return err
	}
	var bases []string
	conn.Raw(fmt.Sprintf("SELECT \"datname\" FROM pg_database WHERE datname = '%s'", dbName)).Scan(&bases)
	if len(bases) == 0 {
		_ = conn.Exec("CREATE DATABASE " + dbName)
	}

	dbUri = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	conn, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
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
	conn, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: false})
	if err != nil {
		return err
	}
	db = conn
	err = db.AutoMigrate(modelsForMigration...)
	if err != nil {
		return err
	}
	return nil
}
