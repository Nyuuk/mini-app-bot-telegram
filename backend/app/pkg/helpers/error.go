package helpers

import (
	"strings"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e Error) Error() string {
	return e.Message
}

func IsDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

func IsNullConstraintError(err error) bool {
	return strings.Contains(err.Error(), "null value in column")
}

func IsNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "record not found")
}

func ErrorClient(message string, code int, data any) Error {
	err := Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
	return err
}

func ErrorInternalServer(data any) Error {
	return Error{
		Code:    500,
		Message: "Internal server error",
		Data:    data,
	}
}
