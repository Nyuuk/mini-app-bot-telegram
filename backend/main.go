package main

import (
	"log"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/controllers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Initializing application...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	log.Println("Connecting to database...")
	if err := database.PGOpen(); err != nil {
		log.Fatal("Error connecting to database")
	}

	// run auto migration
	if err := database.ClientPostgres.AutoMigrate(
		&entities.User{},
		&entities.APIKey{},
		&entities.TelegramUser{},
	); err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	log.Println("Starting server...")
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	authController := controllers.AuthController{}
	userController := controllers.UserController{}

	auth := app.Group("/v1/auth").Name("auth")
	auth.Get("/check", authController.CheckUserApiKey)

	user := app.Group("/v1/user").Name("user")
	user.Get("/:id", userController.GetUserById)
	user.Post("/", userController.CreateUser)

	app.Listen(":3000")
}
