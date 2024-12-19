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
	JsonSyntaxError        *json.SyntaxError
	UnmarshalTypeError     *json.UnmarshalTypeError
	InvalidUnmarshalError  *json.InvalidUnmarshalError
	MaxBytesError          *http.MaxBytesError
)

// Envelope wraps response data for consistent JSON output.
// It uses a map with string keys and any values to provide flexibility
// in the response structure while maintaining JSON compatibility.
type Envelope map[string]any

// Write is a helper that writes a JSON response with the given status code and response data.
// It automatically handles JSON formatting and sets appropriate headers.
// It returns either an error or nil.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to
//   - status: HTTP status code to set in the response
//   - data: The data to be JSON encoded and written
//   - headers: Optional HTTP headers to include in the response (can be nil)
//
// Returns an error if JSON marshaling fails.
//
// Example:
//
//	data := netio.Envelope{"status": "success"}
//
//	headers := http.Header{}
//	headers.Set("X-Custom", "value")
//	headers["X-Many-Values"] = []string{"value1", "value2", "value3"}
//	headers.Add("X-Allowed-Methods", "GET")
//
//	err := netio.Write(w, http.StatusOK, data, headers)
//	if err != nil {
//	  // handle error
//	}
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

func Read(w http.ResponseWriter, r *http.Request, dst any) error {
	// TODO: make this value configurable
	var max int64 = 1_048_576

	// set maximum bytes to receive to prevent/mitigate DOS on API
	http.MaxBytesReader(w, r.Body, int64(max))

	// configure decoder settings
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// decode request body to destination (dst any)
	err := dec.Decode(dst)
	if err != nil {
		return err
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
