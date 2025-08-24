package database

import (
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ClientPostgres *gorm.DB

func PGOpen() error {
	databaseUrl := helpers.GetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres")
	dbClient, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		return err
	}

	ClientPostgres = dbClient

	sqlDB, err := ClientPostgres.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	return nil
}
