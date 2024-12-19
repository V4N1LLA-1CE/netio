// Package netio provides utilities that help with handling json write and read operations
// from HTTP requests and responses.
package netio

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	// ErrNetioMarshalFailure is returned when json marshaling of the response data fails.
	ErrNetioMarshalFailure = errors.New("error marshalling data")
)

// Envelope wraps response data for consistent JSON output.
// It uses a map with string keys and any values to provide flexibility
// in the response structure while maintaining JSON compatibility.
type Envelope map[string]any

// Write is a helper that writes a JSON response with the given status code and response data.
// It automatically handles JSON formatting and sets appropriate headers.
// It returns either an error or nil.
func Write(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

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

	w.Write(json)

	return nil
}

// Read is a helper that reads a JSON response from the request into a dst struct.
// If request size is too big, this helper returns an error.
// Currently, max request size is hardcoded to 1MB for the request body
// It returns either an error or nil.
func Read(w http.ResponseWriter, r *http.Request, dst any) error {
	// TODO: make this value configurable
	var max int64 = 1_048_576

	// set maximum bytes to receive to prevent/mitigate DOS on API
	http.MaxBytesReader(w, r.Body, max)

	// configure decoder settings
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// decode request body to destination (dst any)
	err := dec.Decode(dst)
	if err != nil {
		return fmt.Errorf("read json: %w", err)
	}

	// try decode again into an anonymous dst
	// look for io.EOF. This is to prevent multiple
	// json bodies being used i.e.
	// {"body1": "values"}{"body2": "values"}
	s := &struct{}{}
	err = dec.Decode(s)
	if !errors.Is(err, io.EOF) {
		return errors.New("body can only contain a single json value")
	}

	return nil
}
