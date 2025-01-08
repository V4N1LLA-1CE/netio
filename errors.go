package netio

import (
	"errors"
	"net/http"
	"time"
)

var (
	// NetioUnknownErr is returned when no error is provided to the Error function.
	NetioUnknownErr = errors.New("NetioUnknownError: error not provided")
	// NetioFallbackErr is used when the primary error handling fails.
	NetioFallbackErr = errors.New("NetioErrorFallback: something went wrong...")
	// NetioValidationErr is used when validation checks fail.
	NetioValidationErr = errors.New("validation failed")
)

type ErrorResponse struct {
	// Status represents the HTTP status code
	Status int `json:"status"`
	// StatusText contains the HTTP status text (e.g., "Bad Request" for 400)
	StatusText string `json:"status_text"`
	// Message contains a user-friendly error description
	Message string `json:"message"`
	// ValidationErrors holds validation-specific errors when present
	// This works alongside that netio.Validator
	ValidationErrors any `json:"validation,omitempty"`
	// Timestamp indicates when the error occurred
	Timestamp time.Time `json:"timestamp"`
}

// ErrorFallback returns a generic error response envelope used when primary
// error handling fails. It always returns a 500 Internal Server Error.
func ErrorFallback() Envelope {
	return Envelope{
		"error": ErrorResponse{
			Status:     http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Message:    NetioFallbackErr.Error(),
			Timestamp:  time.Now(),
		},
	}
}

// BuildError creates a new ErrorResponse with the given HTTP status code
// and error message. The Error field is automatically populated with the
// corresponding HTTP status text.
func BuildError(status int, message string) ErrorResponse {
	return ErrorResponse{
		Status:     status,
		StatusText: http.StatusText(status),
		Message:    message,
		Timestamp:  time.Now(),
	}
}

// BuildErrorWithValidation creates a new ErrorResponse that includes validation
// errors from a Validator instance. It includes all fields from BuildError plus
// the validation errors map.
func BuildErrorWithValidation(status int, message string, v *Validator) ErrorResponse {
	return ErrorResponse{
		Status:           status,
		StatusText:       http.StatusText(status),
		Message:          message,
		ValidationErrors: v.Errors,
		Timestamp:        time.Now(),
	}
}

// Error writes a standardized error response to an HTTP response writer. It supports
// both simple errors and validation errors, automatically formatting them into a
// consistent JSON structure.
//
// This utility works in conjunction with netio.Validator for input validation:
//
//	v := netio.NewValidator()
//	v.Check(len(name) >= 2, "name", "too short")
//	v.Check(age >= 18, "age", "must be over 18")
//	if !v.Valid() {
//	    netio.Error(w, "error", nil, http.StatusUnprocessableEntity, v)
//	}
//
// The error response will be wrapped in an envelope using the provided key:
//
//	{
//	    "error": {
//	        "status": 422,
//	        "status_text": "Unprocessable Entity",
//	        "message": "validation failed",
//	        "validation": {     // Only present with validator
//	            "name": "too short",
//	            "age": "must be over 18"
//	        },
//	        "timestamp": "2024-01-08T10:00:00Z"
//	    }
//	}
//
// Parameters:
//
//	w - The response writer to output the error to
//	key - The JSON key to wrap the error in (defaults to "error" if empty)
//	e - The error to include (defaults to NetioUnknownErr if nil)
//	code - The HTTP status code to send (e.g., http.StatusUnprocessableEntity)
//	v - Optional validator containing field-specific validation errors
//
// Example usage:
//
//	Simple error:
//	  netio.Error(w, "error", err, http.StatusBadRequest, nil)
//
//	Validation error:
//	  v := netio.NewValidator()
//	  v.Check(len(name) >= 2, "name", "too short")
//	  if !v.Valid() {
//	      netio.Error(w, "error", nil, http.StatusUnprocessableEntity, v)
//	  }
//
// If the error response cannot be written, it falls back to a generic 500 error
// using ErrorFallback().
func Error(w http.ResponseWriter, key string, e error, code int, v *Validator) {
	if e == nil {
		e = NetioUnknownErr
	}
	if key == "" {
		key = "error"
	}
	// build error response
	var res ErrorResponse
	if v != nil {
		res = BuildErrorWithValidation(code, NetioValidationErr.Error(), v)
	} else {
		res = BuildError(code, e.Error())
	}
	// wrap error with envelope
	env := Envelope{key: res}
	if err := Write(w, code, env, nil); err != nil {
		// if failed to write, fallback to writing generic error
		Write(w, http.StatusInternalServerError, ErrorFallback(), nil)
	}
}
