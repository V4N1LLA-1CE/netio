package netio

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	// structs for testing
	type TestData struct {
		Name       string `json:"name"`
		Age        int    `json:"age"`
		Experience bool   `json:"experience"`
	}

	tests := []struct {
		name       string
		status     int
		data       Envelope
		headers    http.Header
		wantErr    bool
		wantStatus int
	}{
		{
			name:       "simple write with empty header",
			status:     http.StatusOK,
			data:       Envelope{"message": "success"},
			headers:    nil,
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name:   "complex write with headers",
			status: http.StatusOK,
			data: Envelope{"user": TestData{
				Name:       "Jack",
				Age:        42,
				Experience: false,
			}},
			headers: http.Header{
				"X-Many-Headers":  []string{"X-Do-Something", "X-Header-2", "X-Header-3"},
				"X-Single-Header": []string{"X-Do-Not-Cache"},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "fail write with invalid data",
			status:     http.StatusInternalServerError,
			data:       Envelope{"data": func() {}},
			headers:    nil,
			wantErr:    true,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := Write(w, test.status, test.data, test.headers)

			hasErr := false
			if err != nil {
				hasErr = true
			}

			if hasErr != test.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, test.wantErr)
			}

			// make sure status is set
			if w.Code != test.wantStatus {
				t.Errorf("Write() code = %v, want %v", w.Code, test.wantStatus)
			}

			// make sure application/json is always set
			if w.Header().Get("Content-Type") != "application/json" {
				t.Error("Write() did not set Content-Type header to application/json")
			}

			// check if header values are the same
			if test.headers != nil {
				for key, values := range test.headers {
					got := w.Header()[key]
					if !slices.Equal(got, values) {
						t.Errorf("Write() headers don't match what was given")
					}
				}
			}

			// check if json encoded by Write can be decoded without errors
			if !test.wantErr {
				var gotData Envelope
				if err := json.NewDecoder(w.Body).Decode(&gotData); err != nil {
					t.Errorf("Write() invalid JSON response")
				}
			}
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "valid json",
			body:    `{"name": "test", "age": 30}`,
			wantErr: false,
		},
		{
			name:    "invalid json",
			body:    `{"name": "test"`,
			wantErr: true,
		},
		{
			name:    "multiple json objects",
			body:    `{"name": "test"}{"age": 30}`,
			wantErr: true,
		},
		{
			name:    "empty body",
			body:    "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// create request with test body
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			w := httptest.NewRecorder()

			// create destination
			var dst struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}

			err := Read(w, r, &dst)
			if (err != nil) != test.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
