package netio

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrNetioMarshalFailure = errors.New("error marshalling data")
)

type Envelope map[string]any

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
