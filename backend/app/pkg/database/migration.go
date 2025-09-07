package database

import (
	"log"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/joho/godotenv"
)


func RunMigration() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := PGOpen(); err != nil {
		log.Fatal("Error connecting to database")
		return
	}
	err := ClientPostgres.AutoMigrate(
		&entities.User{},
		&entities.APIKey{},
		&entities.TelegramUser{},
		&entities.Overtime{},
		&entities.LogRequest{},
	)
	if err != nil {
		log.Fatal("Error migrating database: ", err)
		return
	}
	log.Println("Migration completed")
}