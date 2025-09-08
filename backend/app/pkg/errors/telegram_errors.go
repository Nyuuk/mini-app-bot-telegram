package errors

// Telegram-specific error types
const (
	// Telegram linking errors
	ErrTelegramAlreadyLinked     = "TELEGRAM_ALREADY_LINKED"
	ErrTelegramInvalidID         = "TELEGRAM_INVALID_ID"
	ErrTelegramUserNotFound      = "TELEGRAM_USER_NOT_FOUND"
	ErrTelegramUnauthorized      = "TELEGRAM_UNAUTHORIZED"
	ErrTelegramDuplicateLink     = "TELEGRAM_DUPLICATE_LINK"
	
	// Telegram validation errors
	ErrTelegramInvalidUsername   = "TELEGRAM_INVALID_USERNAME"
	ErrTelegramInvalidFirstName  = "TELEGRAM_INVALID_FIRST_NAME"
	ErrTelegramInvalidLastName   = "TELEGRAM_INVALID_LAST_NAME"
	
	// Telegram API errors
	ErrTelegramAPIError          = "TELEGRAM_API_ERROR"
	ErrTelegramRateLimit         = "TELEGRAM_RATE_LIMIT"
	
	// Telegram parsing errors
	ErrTelegramIDParseError      = "TELEGRAM_ID_PARSE_ERROR"
)

// TelegramError represents a Telegram-specific error
type TelegramError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func (e TelegramError) Error() string {
	return e.Message
}

// NewTelegramError creates a new Telegram error
func NewTelegramError(code, message string, details any) TelegramError {
	return TelegramError{
		Code:    code,
		Message: message,
		Details: details,
	}
}