package models

import (
	"fmt"
	"net/http"

	"github.com/nbrglm/auth-platform/opts"
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

// WithUI sets the UI parameters for the error if it's an ErrorResponse instance.
// It only sets the RetryUrl, RedirectUrl, and RetryButtonText fields IF they are not already set.
func (e *ErrorResponse) WithUI(retryUrl, redirectUrl, retryButtonText string) *ErrorResponse {
	if e == nil {
		return nil
	}
	if e.UIParams == nil {
		e.UIParams = &ErrorUIParams{
			RetryUrl:        retryUrl,
			RedirectUrl:     redirectUrl,
			RetryButtonText: retryButtonText,
		}
	}
	return e
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
	// Code is an error code that can be used for programmatic handling of errors.
	Code int `json:"code"`

	// UnderlyingError is an optional field that can hold the original error
	// that caused this error response. It is not serialized to JSON.
	// This field is useful for logging and debugging purposes.
	UnderlyingError error `json:"-"`

	UIParams *ErrorUIParams `json:"-"` // UIParams holds additional parameters for UI handling, such as redirect URLs and button text.
}

// ErrorUIParams holds additional parameters for UI handling, such as redirect URLs and button text.
type ErrorUIParams struct {
	// RetryUrl is an internal field used for the UI to redirect the user to a retry page.
	RetryUrl string `json:"-"`

	// RedirectUrl is an internal field used for the UI to redirect the user to a specific page.
	RedirectUrl string `json:"-"`

	// RetryButtonText is an internal field used for the UI to display a retry button with a specific text.
	RetryButtonText string `json:"-"`
}

func (e *ErrorResponse) Error() string {
	if opts.Debug {
		return fmt.Sprintf("Error %d: %s (Debug: %s)", e.Code, e.Message, e.DebugMessage)
	}
	return e.Message
}
