package netio

import "net/http"

// ErrorFallback is a generic error response used when the primary error writing fails.
// It returns a simple "internal server error" message in a JSON envelope.
var ErrorFallback = Envelope{"error": "internal server error"}

// Error writes a JSON error response to the provided http.ResponseWriter.
// It wraps the error in an envelope with an "error_response" key and writes it with
// the specified HTTP status code. If writing fails, it falls back to a generic error
// response with a 500 status code.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to
//   - error: The error value to send. Can be:
//   - string: {"error_response": "error message"}
//   - error: {"error_response": err.Error()}
//   - struct: {"error_response": {"field": "value", ...}}
//   - code: The HTTP status code to use
//
// Example usage:
//
//	// String error
//	Error(w, "invalid input", http.StatusBadRequest)
//
//	// Struct error
//	Error(w, ValidationError{
//	    Field: "email",
//	    Message: "invalid format",
//	}, http.StatusBadRequest)
func Error(w http.ResponseWriter, error any, code int) {
	// wrap error with envelope
	e := Envelope{"error_response": error}
	if err := Write(w, code, e, nil); err != nil {
		// if failed to write, fallback to writing generic error
		Write(w, http.StatusInternalServerError, ErrorFallback, nil)
	}
}
