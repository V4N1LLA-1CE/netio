// Package netio provides utilities that help with handling json write and read operations
// from HTTP requests and responses.
package netio

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	// ErrNetioMarshalFailure is returned when json marshaling of the response data fails.
	ErrNetioMarshalFailure = errors.New("error marshalling data")
)

// An Envelope wraps response data for consistent JSON output.
type Envelope map[string]any

// Write is a helper that writes a JSON response with the given status code and response data.
// It automatically handles JSON formatting and sets appropriate headers.
// It returns either an error or nil.
func Write(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return ErrNetioMarshalFailure
	}

	// formatting for terminal i.e. curl responses
	json = append(json, '\n')

	// go through headers map and apply headers
	for key, values := range headers {
		w.Header()[key] = values
	}

	// necessary headers for json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)

	return nil
}
