package main

import (
	"log"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/controllers"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/middlewares"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/database"
	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Mini App Bot Telegram API
// @version 1.0
// @description API for Mini App Bot Telegram Backend System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
	// if err := database.ClientPostgres.AutoMigrate(
	// 	&entities.User{},
	// 	&entities.APIKey{},
	// 	&entities.TelegramUser{},
	// 	&entities.Overtime{},
	// ); err != nil {
	// 	log.Fatal("Error migrating database: ", err)
	// }

	log.Println("Starting server...")
	app := fiber.New()

	// Apply logging middlewares
	app.Use(middlewares.LoggingMiddleware())
	app.Use(middlewares.DatabaseLoggingMiddleware())
	app.Use(middlewares.PerformanceLoggingMiddleware())
	app.Use(middlewares.SecurityLoggingMiddleware())

	// Swagger UI route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Health check endpoint
	// @Summary Health Check
	// @Description Check if the API is running
	// @Tags Health
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	authController := controllers.AuthController{}
	userController := controllers.UserController{}
	telegramController := controllers.TelegramController{}
	overtimeController := controllers.OvertimeController{}

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

	// Overtime routes
	overtime := protected.Group("/overtime").Name("overtime")
	overtime.Post("/", overtimeController.CreateNewRecordOvertime)                              // Create new overtime record
	overtime.Get("/telegram/:telegram_id", overtimeController.GetAllRecordOvertimeByTelegramID) // Get all overtime records by telegram ID
	overtime.Post("/by-date", overtimeController.GetRecordByDateByTelegramID)                   // Get overtime record by specific date
	overtime.Post("/between-dates", overtimeController.GetRecordBetweenDateByTelegramId)        // Get overtime records between dates
	overtime.Get("/:id", overtimeController.GetRecordByID)                                      // Get overtime record by ID
	overtime.Put("/:id", overtimeController.UpdateRecordOvertime)                               // Update overtime record
	overtime.Delete("/:id", overtimeController.DeleteRecordOvertime)                            // Delete overtime record

	// API Key routes
	// apikey := protected.Group("/apikey").Name("apikey")
	// apikey.Get("/", authController.GetUserApiKeys)     // Get semua API key user
	// apikey.Post("/", authController.CreateApiKey)      // Create API key baru
	// apikey.Delete("/:id", authController.DeleteApiKey) // Delete API key

	app.Listen(":3000")
}
