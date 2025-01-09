package main

import (
	"github.com/V4N1LLA-1CE/netio"
	"log"
	"net/http"
	"regexp"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", registerHandler)
	log.Printf("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// read JSON input
	var input struct {
		Email string `json:"email"`
		Age   int    `json:"age"`
	}

	// read request body into input struct
	if err := netio.Read(w, r, &input); err != nil {
		netio.Error(w, "error", http.StatusBadRequest, nil)
		return
	}

	// validate input
	v := netio.NewValidator()
	emailRx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	v.Check(netio.Matches(input.Email, emailRx), "email", "invalid email format")
	v.Check(input.Age >= 18, "age", "must be over 18")

	if !v.Valid() {
		netio.Error(w, "error", http.StatusUnprocessableEntity, v)
		return
	}

	// success response
	response := netio.Envelope{
		"success": map[string]any{
			"message": "User registered successfully",
			"user":    input,
		},
	}
	if err := netio.Write(w, http.StatusCreated, response, nil); err != nil {
		netio.Error(w, "error", http.StatusInternalServerError, nil)
	}
}
