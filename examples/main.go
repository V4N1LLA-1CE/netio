package main

import (
	"github.com/V4N1LLA-1CE/netio"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/example", exampleHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Printf("Server starting on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func exampleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string         `json:"username"`
		UserData map[string]any `json:"user_data"`
	}
	// read data into input struct
	// netio.Read() will
	err := netio.Read(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// write response
	//
	// headers can be set to  nil for netio.Write()
	// if no further configuration is needed
	headers := http.Header{
		"X-Some-Header": []string{"Value-1", "Value-2"},
		"X-API-Version": []string{"1.0"},
	}
	headers.Add("X-New-Header", "New-Header-Value")
	// netio.Write() will automatically set json headers (can be overriden by custom header)
	//
	// any type of data to be written must be wrapped by netio.Envelope
	// which is just a map[string]any
	// this ensures good JSON response structures
	err = netio.Write(w, http.StatusOK, netio.Envelope{"example response": input}, headers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
