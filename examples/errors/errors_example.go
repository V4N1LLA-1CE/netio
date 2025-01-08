package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/V4N1LLA-1CE/netio"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/errors", errorsHandler)

	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func errorsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Query().Get("type") {
	case "simple":
		// Simple error
		err := errors.New("something went wrong")
		netio.Error(w, "error", err, http.StatusBadRequest, nil)

	case "validation":
		// Error with validation
		v := netio.NewValidator()
		// Simulate false conditions
		v.Check(false, "email", "invalid email format")
		v.Check(false, "age", "must be over 18")

		err := errors.New("validation failed")
		netio.Error(w, "error", err, http.StatusUnprocessableEntity, v)

	default:
		// Fallback error
		netio.Write(w, http.StatusInternalServerError, netio.ErrorFallback(), nil)
	}
}
