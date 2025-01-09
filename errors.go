package netio

import (
	"net/http"
	"time"
)

// ErrorResponse represents a standardized error response structure for HTTP APIs.
// It includes the status code, message, optional validation errors, and timestamp
// of when the error occurred.
type ErrorResponse struct {
	// Status represents the HTTP status code
	Status int `json:"status"`
	// Message contains the HTTP status text (e.g., "Bad Request" for 400)
	Message string `json:"message"`
	// ValidationErrors holds validation-specific errors when present.
	// This field works in conjunction with netio.Validator to provide
	// detailed validation feedback to API clients.
	ValidationErrors any `json:"validation,omitempty"`
	// Timestamp indicates when the error occurred
	Timestamp time.Time `json:"timestamp"`
}

// ErrorFallback returns a generic error response envelope used when primary
// error handling fails. It always returns a 500 Internal Server Error wrapped
// in an Envelope.
func ErrorFallback() Envelope {
	return Envelope{
		"error": ErrorResponse{
			Status:    http.StatusInternalServerError,
			Message:   http.StatusText(http.StatusInternalServerError),
			Timestamp: time.Now(),
		},
	}
}

// BuildError creates a new ErrorResponse with the specified HTTP status code.
// The message is automatically set to the standard HTTP status text for the given code.
func BuildError(status int) ErrorResponse {
	return ErrorResponse{
		Status:    status,
		Message:   http.StatusText(status),
		Timestamp: time.Now(),
	}
}

// BuildErrorWithValidation creates a new ErrorResponse that includes validation errors.
// It combines the HTTP status code with validation errors from a Validator instance.
func BuildErrorWithValidation(status int, v *Validator) ErrorResponse {
	return ErrorResponse{
		Status:           status,
		Message:          http.StatusText(status),
		ValidationErrors: v.Errors,
		Timestamp:        time.Now(),
	}
}

// Error writes a JSON error response to the provided http.ResponseWriter.
// It handles various error scenarios and provides a consistent error structure
// across your API endpoints.
//
// The error response is wrapped in an envelope and includes:
//   - HTTP status code and message
//   - Timestamp of when the error occurred
//   - Optional validation errors from netio.Validator
//
// If writing the response fails, it falls back to a generic 500 Internal Server Error.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to
//   - key: The JSON key for wrapping the error in the response envelope (defaults to "error" if empty)
//   - code: The HTTP status code (automatically corrected to 500 if invalid)
//   - v: Optional *netio.Validator containing validation errors (may be nil)
//
// Basic usage example:
//
//	// Return a simple 404 error
//	netio.Error(w, "error", http.StatusNotFound, nil)
//
// Example with validation:
//
//	v := netio.NewValidator()
//
//	// Validate user input
//	v.Check(len(user.Password) >= 8, "password", "must be at least 8 characters")
//	v.Check(netio.Matches(user.Email, emailRX), "email", "must be a valid email")
//
//	if !v.Valid() {
//	    // Returns a 400 Bad Request with validation details
//	    netio.Error(w, "error", http.StatusBadRequest, v)
//	    return
//	}
//
// The JSON response format for validation errors:
//
//	{
//	    "error": {
//	        "status": 400,
//	        "message": "Bad Request",
//	        "validation": {
//	            "password": "must be at least 8 characters",
//	            "email": "must be a valid email"
//	        },
//	        "timestamp": "2024-01-09T12:00:00Z"
//	    }
//	}
//
// The JSON response format for non-validation errors:
//
//	{
//	    "error": {
//	        "status": 404,
//	        "message": "Not Found",
//	        "timestamp": "2024-01-09T12:00:00Z"
//	    }
//	}
func Error(w http.ResponseWriter, key string, code int, v *Validator) {
	// handle invalid code and empty key
	if code < 100 || code > 599 {
		code = http.StatusInternalServerError
	}
	if key == "" {
		key = "error"
	}
	// build error response
	var res ErrorResponse
	if v != nil {
		res = BuildErrorWithValidation(code, v)
	} else {
		res = BuildError(code)
	}
	// wrap error with envelope
	env := Envelope{key: res}
	if err := Write(w, code, env, nil); err != nil {
		// if failed to write, fallback to writing generic error
		Write(w, http.StatusInternalServerError, ErrorFallback(), nil)
	}
}
