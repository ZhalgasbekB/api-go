package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorPage(t *testing.T) {
	app := &API{}

	tests := []struct {
		name           string
		statusCode     int
		errorMessage   string
		expectedOutput string
	}{
		{
			name:           "404 Not Found",
			statusCode:     http.StatusNotFound,
			errorMessage:   "Page not found",
			expectedOutput: "Page not found",
		},
		{
			name:           "500 Internal Server Error",
			statusCode:     http.StatusInternalServerError,
			errorMessage:   "Internal server error",
			expectedOutput: "Internal server error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			app.Error(w, test.statusCode, test.errorMessage)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != test.statusCode {
				t.Errorf("Expected status code %d, got %d", test.statusCode, resp.StatusCode)
			}

			if !bytes.Contains(body, []byte(test.expectedOutput)) {
				t.Errorf("Expected response body to contain %s, got %s", test.expectedOutput, string(body))
			}
		})
	}
}
