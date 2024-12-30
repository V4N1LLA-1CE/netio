// Package netio provides lightweight utilities that simplify common
// webserver development tasks in Go. This package aims to reduce boilerplate
// code for simple components so the developer can focus on application logic.
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
	// ErrMultipleJsonBodies is returned when Read detects more than one body i.e. {}{}
	ErrMultipleJsonBodies = errors.New("body can only contain a single json value")
)

// Envelope represents a wrapper for HTTP response data in JSON format.
// It uses a map with string keys and interface{} values to provide
// flexibility in the response structure while maintaining JSON compatibility.
//
// Example:
//
//	type metadata struct {
//	  ...
//	}
//
//	m := metadata{...}
//
//	// good practice to send JSON with envelope wrapper
//	responseData := netio.Envelope{
//	    "metadata": m,
//	}
type Envelope map[string]any

// Write sends a JSON response with the given status code and response data.
// It handles JSON formatting, sets appropriate headers, and provides pretty-printing
// for better CLI tool readability.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to
//   - status: HTTP status code to send
//   - data: The Envelope containing response data to be JSON encoded
//   - headers: Additional HTTP headers to include in the response
//
// Returns an error if JSON marshaling fails, nil otherwise.
//
// Example:
//
//	env := netio.Envelope{"users": users}
//	headers := http.Header{"X-Custom": []string{"value"}}
//	err := netio.Write(w, http.StatusOK, env, headers)
func Write(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {
	// header good practices (OWASP)
	// see more at https://cheatsheetseries.owasp.org/cheatsheets/HTTP_Headers_Cheat_Sheet.html
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

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

// Read decodes a JSON request body into the provided destination struct.
// It enforces a maximum request size of 1MB and validates that only a single
// JSON object is present in the request body.
//
// Parameters:
//   - w: The http.ResponseWriter (used for MaxBytesReader)
//   - r: The *http.Request containing the JSON body
//   - dst: Non-nil pointer to the destination struct where the JSON will be decoded
//
// Example:
//
//	var input struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	if err := netio.Read(w, r, &input); err != nil {
//	    // Handle error...
//	}
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
		return fmt.Errorf("netio.Read(): %w", err)
	}

	// try decode again into an anonymous dst
	// look for io.EOF. This is to prevent multiple
	// json bodies being used i.e.
	// {"body1": "values"}{"body2": "values"}
	s := &struct{}{}
	err = dec.Decode(s)
	if !errors.Is(err, io.EOF) {
		return ErrMultipleJsonBodies
	}

	return nil
}
