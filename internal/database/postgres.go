package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	// Load environment variables
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Bangkok",
		host, user, password, dbName, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to the database with GORM")
	DB = db
	return db, nil
}

func Migrate(db *gorm.DB, dst ...interface{}) error {
	err := db.AutoMigrate(dst...)
	if err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
