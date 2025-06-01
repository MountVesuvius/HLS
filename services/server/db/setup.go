package db

import (
	"fmt"
	"log"
	"os"
	"server/server/dto"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
    user := os.Getenv("DB_USER")
    name := os.Getenv("DB_NAME")
    password := os.Getenv("DB_PASSWORD")
    host := os.Getenv("DB_HOST")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, name)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database.", err)
    }
    return db
}

// will move to migration scripts once i have a clear idea of how this will look
func Sync(db *gorm.DB) {
    err := db.AutoMigrate(&dto.File{})
    if err != nil {
        log.Fatal("Failed to sync database.", err)
    }
}
