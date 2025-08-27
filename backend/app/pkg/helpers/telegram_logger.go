package helpers

import (
	"time"
)

// LogTelegramMessage logs incoming Telegram messages
func LogTelegramMessage(chatID int64, userID int64, messageType string, text string, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_message").
		Int64("chat_id", chatID).
		Int64("user_id", userID).
		Str("message_type", messageType).
		Str("text", text)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Message Received")
}

// LogTelegramCommand logs Telegram bot commands
func LogTelegramCommand(chatID int64, userID int64, command string, arguments string, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_command").
		Int64("chat_id", chatID).
		Int64("user_id", userID).
		Str("command", command).
		Str("arguments", arguments)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Command Executed")
}

// LogTelegramCallback logs Telegram callback queries
func LogTelegramCallback(chatID int64, userID int64, callbackData string, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_callback").
		Int64("chat_id", chatID).
		Int64("user_id", userID).
		Str("callback_data", callbackData)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Callback Query")
}

// LogTelegramBotAction logs bot actions (sending messages, etc.)
func LogTelegramBotAction(action string, chatID int64, userID int64, success bool, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_bot_action").
		Str("action", action).
		Int64("chat_id", chatID).
		Int64("user_id", userID).
		Bool("success", success)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	if success {
		logEvent.Msg("Telegram Bot Action Success")
	} else {
		logEvent.Msg("Telegram Bot Action Failed")
	}
}

// LogTelegramUserRegistration logs user registration events
func LogTelegramUserRegistration(chatID int64, userID int64, username string, firstName string, lastName string, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_user_registration").
		Int64("chat_id", chatID).
		Int64("user_id", userID).
		Str("username", username).
		Str("first_name", firstName).
		Str("last_name", lastName)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram User Registration")
}

// LogTelegramError logs Telegram-related errors
func LogTelegramError(err error, context string, chatID int64, userID int64, details map[string]interface{}) {
	logEvent := Logger.Error().
		Err(err).
		Str("type", "telegram_error").
		Str("context", context).
		Int64("chat_id", chatID).
		Int64("user_id", userID)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Error")
}

// LogTelegramWebhook logs webhook events
func LogTelegramWebhook(eventType string, updateID int64, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_webhook").
		Str("event_type", eventType).
		Int64("update_id", updateID)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Webhook Event")
}

// LogTelegramMiniApp logs mini app interactions
func LogTelegramMiniApp(chatID int64, userID int64, appName string, action string, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_mini_app").
		Int64("chat_id", chatID).
		Int64("user_id", userID).
		Str("app_name", appName).
		Str("action", action)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Mini App Interaction")
}

// LogTelegramPerformance logs performance metrics for Telegram operations
func LogTelegramPerformance(operation string, duration time.Duration, chatID int64, userID int64, details map[string]interface{}) {
	logEvent := Logger.Info().
		Str("type", "telegram_performance").
		Str("operation", operation).
		Dur("duration", duration).
		Int64("chat_id", chatID).
		Int64("user_id", userID)

	// Add additional details
	for key, value := range details {
		logEvent = logEvent.Interface(key, value)
	}

	logEvent.Msg("Telegram Performance Metric")
}
