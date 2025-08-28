package main

import (
	"log"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/controllers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/entities"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/middlewares"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Initializing application...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	helpers.InitLogger()

	log.Println("Connecting to database...")
	if err := database.PGOpen(); err != nil {
		log.Fatal("Error connecting to database")
	}

	// run auto migration
	if err := database.ClientPostgres.AutoMigrate(
		&entities.User{},
		&entities.APIKey{},
		&entities.TelegramUser{},
		&entities.Overtime{},
	); err != nil {
		log.Fatal("Error migrating database: ", err)
	}

	log.Println("Starting server...")
	app := fiber.New()

	// Apply logging middlewares
	app.Use(middlewares.LoggingMiddleware())
	app.Use(middlewares.DatabaseLoggingMiddleware())
	app.Use(middlewares.PerformanceLoggingMiddleware())
	app.Use(middlewares.SecurityLoggingMiddleware())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	authController := controllers.AuthController{}
	userController := controllers.UserController{}
	telegramController := controllers.TelegramController{}

	// Public routes (tidak perlu auth)
	auth := app.Group("/v1/auth").Name("auth")
	auth.Post("/login", authController.Login)         // Login untuk dapat JWT
	auth.Post("/register", userController.CreateUser) // Register user baru

	// Protected routes (perlu auth via API Key atau JWT)
	protected := app.Group("/v1", middlewares.AuthMiddleware()).Name("protected")

	// User routes
	user := protected.Group("/user").Name("user")
	user.Get("/detail-me", userController.GetDetailMe)           // Get all users (admin only)
	user.Post("/", userController.CreateUser)                    // Create user (admin only)
	user.Get("/", userController.GetAllUsers)                    // Get all users (admin only)
	user.Get("/api-key", userController.GetApiKeyFromUserActive) // Get API key from user active
	user.Post("/api-key", userController.CreateApiKey)           // Create API key baru

	user.Get("/:id", userController.GetUserById)       // Get user by ID (admin only)
	user.Delete("/:id", userController.DeleteUserById) // Delete user by ID (admin only)

	telegram := protected.Group("/telegram").Name("telegram")
	telegram.Post("/", telegramController.CreateNewUserForNowUserActive) // Create new user for now user active
	telegram.Get("/", telegramController.FindByUserID)                   // Get all user telegram
	telegram.Get("/:id", telegramController.FindByTelegramID)            // Get user telegram by ID
	telegram.Delete("/:id", telegramController.DeleteByTelegramID)       // Delete user telegram by ID
	telegram.Put("/:id", telegramController.UpdateByTelegramID)          // Update user telegram by ID

	// API Key routes
	// apikey := protected.Group("/apikey").Name("apikey")
	// apikey.Get("/", authController.GetUserApiKeys)     // Get semua API key user
	// apikey.Post("/", authController.CreateApiKey)      // Create API key baru
	// apikey.Delete("/:id", authController.DeleteApiKey) // Delete API key

	app.Listen(":3000")
}
