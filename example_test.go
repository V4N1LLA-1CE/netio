package netio_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/V4N1LLA-1CE/netio"
)

// ExampleWrite demonstrates how to use the Write function with custom header configurations.
func ExampleWrite() {
	w := httptest.NewRecorder()

	// create custom header
	headers := http.Header{}
	headers.Set("X-Custom-Header", "some-value")
	headers["X-Many-Values"] = []string{"value1", "value2", "value3"}

	headers.Add("X-Allowed-Methods", "GET")
	headers.Add("X-Allowed-Methods", "POST")
	headers.Add("X-Allowed-Methods", "PUT")
	headers.Add("X-Allowed-Methods", "OPTIONS")

	data := netio.Envelope{
		"status": "success",
		"user": map[string]any{
			"id":   1,
			"name": "Test User",
		},
	}

	err := netio.Write(w, http.StatusOK, data, headers)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// print headers
	fmt.Println("Single header:", w.Header().Get("X-Custom-Header"))
	fmt.Println("Multiple headers (first method):", w.Header()["X-Many-Values"])
	fmt.Println("Multiple headers (second method):", w.Header()["X-Allowed-Methods"])
	// Output:
	// Single header: some-value
	// Multiple headers (first method): [value1 value2 value3]
	// Multiple headers (second method): [GET POST PUT OPTIONS]
}
