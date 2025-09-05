package models

import (
	"fmt"
	"net/http"

	"github.com/nbrglm/nexeres/opts"
)

// NewErrorResponse creates a new ErrorResponse instance.
//
// It takes a user-friendly message, a debug message for developers, and an error code.
//
// The debug message is intended for internal use and should not be exposed to end users in production environments.
// For doing that, use the Filter method on the ErrorResponse instance when passing it to gin or the client.
//
// Also, the "Debug" Message is used to record things in span context and/or logs for observability purposes,
// so do not include sensitive information in the debug message.
//
// If underlying error is not nil, it will be logged at ERROR level for debugging purposes,
// If underlying error is nil, no error will be logged.
//
// Note: To specify the `RetryUrl`, `RedirectUrl`, and `RetryButtonText` fields, use the `errResponse.WithUI()` function.
func NewErrorResponse(message, debug string, code int, underlying error) *ErrorResponse {
	return &ErrorResponse{
		Message:         message,
		DebugMessage:    debug,
		Code:            code,
		UnderlyingError: underlying,
	}
}

const GenericErrorMessage = "An error occurred while processing your request. Please try again later."

// NewGenericErrorResponse creates a new ErrorResponse instance with a generic error message.
func NewGenericErrorResponse(opName string, err error) *ErrorResponse {
	return &ErrorResponse{
		Message:         GenericErrorMessage,
		DebugMessage:    "Failed to handle " + opName,
		Code:            http.StatusInternalServerError,
		UnderlyingError: err,
	}
}

// Filter filters the ErrorResponse based on the debug mode.
// If debug mode is enabled, it returns the full error response including the debug message.
// If debug mode is not enabled, it returns a filtered error response without the debug message.
// This is useful for controlling the visibility of debug information in a production environment.
func (e *ErrorResponse) Filter() *ErrorResponse {
	if opts.Debug || e.DebugMessage == "" {
		return e
	}

	// Otherwise, return a filtered error response without the debug message
	return &ErrorResponse{
		Message:      e.Message,
		DebugMessage: "",
		Code:         e.Code,
	}
}

type ErrorResponse struct {
	// Message is a user-friendly message that can be displayed to the end user.
	Message string `json:"message"`
	// DebugMessage is a technical message that can be used for debugging.
	DebugMessage string `json:"debug"`

	// HTTP Status code associated with this error.
	// It is not serialized to JSON. This field is useful for setting the HTTP status code in the response.
	Code int `json:"-"`

	// UnderlyingError is an optional field that can hold the original error
	// that caused this error response. It is not serialized to JSON.
	// This field is useful for logging and debugging purposes.
	UnderlyingError error `json:"-"`
}

func (e *ErrorResponse) Error() string {
	if opts.Debug {
		return fmt.Sprintf("Error %d: %s (Debug: %s)", e.Code, e.Message, e.DebugMessage)
	}
	return e.Message
}
