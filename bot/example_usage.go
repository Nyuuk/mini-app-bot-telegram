package main

import (
	"time"

	"github.com/Nyuuk/mini-app-bot-telegram/backend/app/pkg/helpers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Example of how to use the logging system in a Telegram bot
func handleMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	start := time.Now()

	// Log incoming message
	helpers.LogTelegramMessage(
		update.Message.Chat.ID,
		update.Message.From.ID,
		"text",
		update.Message.Text,
		map[string]interface{}{
			"message_id": update.Message.MessageID,
			"date":       update.Message.Date,
			"username":   update.Message.From.UserName,
			"first_name": update.Message.From.FirstName,
			"last_name":  update.Message.From.LastName,
		},
	)

	// Handle commands
	if update.Message.IsCommand() {
		helpers.LogTelegramCommand(
			update.Message.Chat.ID,
			update.Message.From.ID,
			update.Message.Command(),
			update.Message.CommandArguments(),
			map[string]interface{}{
				"message_id": update.Message.MessageID,
			},
		)

		// Handle specific commands
		switch update.Message.Command() {
		case "start":
			handleStartCommand(bot, update)
		case "help":
			handleHelpCommand(bot, update)
		default:
			handleUnknownCommand(bot, update)
		}
	} else {
		// Handle regular messages
		handleRegularMessage(bot, update)
	}

	// Log performance
	helpers.LogTelegramPerformance(
		"message_handling",
		time.Since(start),
		update.Message.Chat.ID,
		update.Message.From.ID,
		map[string]interface{}{
			"message_type": "text",
		},
	)
}

func handleStartCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	start := time.Now()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! I'm your bot assistant.")
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)

	// Log bot action
	helpers.LogTelegramBotAction(
		"send_message",
		update.Message.Chat.ID,
		update.Message.From.ID,
		err == nil,
		map[string]interface{}{
			"command":  "start",
			"duration": time.Since(start),
		},
	)

	if err != nil {
		helpers.LogTelegramError(
			err,
			"send_start_message",
			update.Message.Chat.ID,
			update.Message.From.ID,
			map[string]interface{}{
				"command": "start",
			},
		)
	}
}

func handleHelpCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	start := time.Now()

	helpText := `Available commands:
/start - Start the bot
/help - Show this help message
/register - Register your account
/profile - View your profile`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpText)
	msg.ParseMode = "HTML"

	_, err := bot.Send(msg)

	// Log bot action
	helpers.LogTelegramBotAction(
		"send_help_message",
		update.Message.Chat.ID,
		update.Message.From.ID,
		err == nil,
		map[string]interface{}{
			"command":  "help",
			"duration": time.Since(start),
		},
	)
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	helpers.LogTelegramCommand(
		update.Message.Chat.ID,
		update.Message.From.ID,
		update.Message.Command(),
		update.Message.CommandArguments(),
		map[string]interface{}{
			"status": "unknown_command",
		},
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command. Use /help to see available commands.")
	bot.Send(msg)
}

func handleRegularMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Log business event
	helpers.LogBusiness(
		"telegram_message_received",
		"", // user ID will be set if user is registered
		map[string]interface{}{
			"chat_id":          update.Message.Chat.ID,
			"telegram_user_id": update.Message.From.ID,
			"message_text":     update.Message.Text,
			"username":         update.Message.From.UserName,
		},
	)

	// Echo the message back
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You said: "+update.Message.Text)
	bot.Send(msg)
}

func handleCallbackQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	start := time.Now()

	// Log callback query
	helpers.LogTelegramCallback(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.From.ID,
		update.CallbackQuery.Data,
		map[string]interface{}{
			"callback_query_id": update.CallbackQuery.ID,
			"message_id":        update.CallbackQuery.Message.MessageID,
		},
	)

	// Handle different callback data
	switch update.CallbackQuery.Data {
	case "register":
		handleRegisterCallback(bot, update)
	case "profile":
		handleProfileCallback(bot, update)
	default:
		handleUnknownCallback(bot, update)
	}

	// Log performance
	helpers.LogTelegramPerformance(
		"callback_handling",
		time.Since(start),
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.From.ID,
		map[string]interface{}{
			"callback_data": update.CallbackQuery.Data,
		},
	)
}

func handleRegisterCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Log user registration attempt
	helpers.LogTelegramUserRegistration(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.From.ID,
		update.CallbackQuery.From.UserName,
		update.CallbackQuery.From.FirstName,
		update.CallbackQuery.From.LastName,
		map[string]interface{}{
			"source": "callback_query",
		},
	)

	// Send registration form or redirect to mini app
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Please register using our mini app.")
	bot.Send(msg)
}

func handleProfileCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Log profile view attempt
	helpers.LogBusiness(
		"profile_view_attempt",
		"", // user ID will be set if user is registered
		map[string]interface{}{
			"chat_id":          update.CallbackQuery.Message.Chat.ID,
			"telegram_user_id": update.CallbackQuery.From.ID,
			"source":           "callback_query",
		},
	)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Your profile information will be shown here.")
	bot.Send(msg)
}

func handleUnknownCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	helpers.LogTelegramCallback(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.From.ID,
		update.CallbackQuery.Data,
		map[string]interface{}{
			"status": "unknown_callback",
		},
	)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Unknown action.")
	bot.Send(msg)
}

func handleMiniAppInteraction(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	// Log mini app interaction
	helpers.LogTelegramMiniApp(
		update.Message.Chat.ID,
		update.Message.From.ID,
		"your_mini_app_name",
		"interaction",
		map[string]interface{}{
			"web_app_data": update.Message.WebAppData,
			"message_id":   update.Message.MessageID,
		},
	)

	// Handle mini app data
	if update.Message.WebAppData != nil {
		// Process web app data
		helpers.LogBusiness(
			"mini_app_data_processed",
			"", // user ID will be set if user is registered
			map[string]interface{}{
				"chat_id":          update.Message.Chat.ID,
				"telegram_user_id": update.Message.From.ID,
				"data":             update.Message.WebAppData.Data,
			},
		)
	}
}
