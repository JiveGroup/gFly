package test

import (
	"encoding/json"
	"gfly/app/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestHttpBuildURL(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		queryParams map[string]string
		expected    string
		expectError bool
	}{
		{
			name:        "Simple URL without query params",
			baseURL:     "https://example.com",
			queryParams: map[string]string{},
			expected:    "https://example.com",
			expectError: false,
		},
		{
			name:    "URL with query params",
			baseURL: "https://example.com",
			queryParams: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			expected:    "https://example.com?param1=value1&param2=value2",
			expectError: false,
		},
		{
			name:    "URL with existing query params",
			baseURL: "https://example.com?existing=value",
			queryParams: map[string]string{
				"param1": "value1",
			},
			expected:    "https://example.com?existing=value&param1=value1",
			expectError: false,
		},
		{
			name:        "Invalid URL",
			baseURL:     "://invalid-url",
			queryParams: map[string]string{},
			expected:    "",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := utils.HttpBuildURL(test.baseURL, test.queryParams)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check if the URL contains all the expected query parameters
				// We don't check for exact string match because the order of query parameters might differ
				if result != test.expected {
					// If the URLs don't match exactly, parse them and compare the components
					// This is a simplified check - in a real test you might want to parse the URLs and compare components
					t.Errorf("Expected URL %s, got %s", test.expected, result)
				}
			}
		})
	}
}

func TestHttpIsSuccessStatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		expected   bool
	}{
		{"Status 200", 200, true},
		{"Status 201", 201, true},
		{"Status 299", 299, true},
		{"Status 300", 300, false},
		{"Status 400", 400, false},
		{"Status 500", 500, false},
		{"Status 199", 199, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.HttpIsSuccessStatusCode(test.statusCode)
			if result != test.expected {
				t.Errorf("Expected %v for status code %d, got %v", test.expected, test.statusCode, result)
			}
		})
	}
}

func TestHttpParseQueryParams(t *testing.T) {
	tests := []struct {
		name        string
		queryString string
		expected    map[string]string
		expectError bool
	}{
		{
			name:        "Empty query string",
			queryString: "",
			expected:    map[string]string{},
			expectError: false,
		},
		{
			name:        "Simple query string",
			queryString: "param1=value1&param2=value2",
			expected: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			expectError: false,
		},
		{
			name:        "Query string with repeated params",
			queryString: "param1=value1&param1=value2",
			expected: map[string]string{
				"param1": "value1", // Only the first value is kept
			},
			expectError: false,
		},
		{
			name:        "Query string with URL encoding",
			queryString: "param1=value%20with%20spaces&param2=special%26chars",
			expected: map[string]string{
				"param1": "value with spaces",
				"param2": "special&chars",
			},
			expectError: false,
		},
		{
			name:        "Invalid query string",
			queryString: "param1=%invalid",
			expected:    nil,
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := utils.HttpParseQueryParams(test.queryString)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(result, test.expected) {
					t.Errorf("Expected %v, got %v", test.expected, result)
				}
			}
		})
	}
}

func TestHttpCreateHTTPClient(t *testing.T) {
	tests := []struct {
		name                string
		timeout             time.Duration
		maxIdleConns        int
		maxIdleConnsPerHost int
		maxConnsPerHost     int
	}{
		{
			name:                "Default client",
			timeout:             10 * time.Second,
			maxIdleConns:        10,
			maxIdleConnsPerHost: 5,
			maxConnsPerHost:     100,
		},
		{
			name:                "Custom client",
			timeout:             30 * time.Second,
			maxIdleConns:        20,
			maxIdleConnsPerHost: 10,
			maxConnsPerHost:     200,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := utils.HttpCreateHTTPClient(test.timeout, test.maxIdleConns, test.maxIdleConnsPerHost, test.maxConnsPerHost)

			// Check timeout
			if client.Timeout != test.timeout {
				t.Errorf("Expected timeout %v, got %v", test.timeout, client.Timeout)
			}

			// Check transport settings
			transport, ok := client.Transport.(*http.Transport)
			if !ok {
				t.Errorf("Expected *http.Transport, got %T", client.Transport)
				return
			}

			if transport.MaxIdleConns != test.maxIdleConns {
				t.Errorf("Expected MaxIdleConns %d, got %d", test.maxIdleConns, transport.MaxIdleConns)
			}

			if transport.MaxIdleConnsPerHost != test.maxIdleConnsPerHost {
				t.Errorf("Expected MaxIdleConnsPerHost %d, got %d", test.maxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
			}

			if transport.MaxConnsPerHost != test.maxConnsPerHost {
				t.Errorf("Expected MaxConnsPerHost %d, got %d", test.maxConnsPerHost, transport.MaxConnsPerHost)
			}
		})
	}
}

// Mock server for testing HTTP requests
func setupMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
	})
	return server
}

func TestHttpGetJSON(t *testing.T) {
	type TestResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	// Setup mock server
	server := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Check custom header if present
		if r.URL.Path == "/custom-header" && r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header to be custom-value, got %s", r.Header.Get("X-Custom-Header"))
		}

		// Return response based on path
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Success",
				Status:  "OK",
			})
		case "/custom-header":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Custom header received",
				Status:  "OK",
			})
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Error",
				Status:  "Error",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	tests := []struct {
		name        string
		url         string
		headers     map[string]string
		expectError bool
		expected    TestResponse
	}{
		{
			name:        "Successful request",
			url:         server.URL + "/success",
			headers:     nil,
			expectError: false,
			expected: TestResponse{
				Message: "Success",
				Status:  "OK",
			},
		},
		{
			name:        "Request with custom headers",
			url:         server.URL + "/custom-header",
			headers:     map[string]string{"X-Custom-Header": "custom-value"},
			expectError: false,
			expected: TestResponse{
				Message: "Custom header received",
				Status:  "OK",
			},
		},
		{
			name:        "Server error",
			url:         server.URL + "/error",
			headers:     nil,
			expectError: false, // The function doesn't return an error for non-2xx status codes
			expected: TestResponse{
				Message: "Error",
				Status:  "Error",
			},
		},
		{
			name:        "Invalid URL",
			url:         "http://invalid-url-that-does-not-exist.example",
			headers:     nil,
			expectError: true,
			expected:    TestResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response TestResponse
			err := utils.HttpGetJSON(test.url, &response, test.headers)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if response.Message != test.expected.Message || response.Status != test.expected.Status {
					t.Errorf("Expected response %+v, got %+v", test.expected, response)
				}
			}
		})
	}
}

func TestHttpPostJSON(t *testing.T) {
	type TestRequest struct {
		Data string `json:"data"`
	}

	type TestResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Echo    string `json:"echo,omitempty"`
	}

	// Setup mock server
	server := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Check custom header if present
		if r.URL.Path == "/custom-header" && r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header to be custom-value, got %s", r.Header.Get("X-Custom-Header"))
		}

		// Parse request body
		var request TestRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return response based on path
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Success",
				Status:  "OK",
				Echo:    request.Data,
			})
		case "/custom-header":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Custom header received",
				Status:  "OK",
				Echo:    request.Data,
			})
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Error",
				Status:  "Error",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	tests := []struct {
		name        string
		url         string
		body        TestRequest
		headers     map[string]string
		expectError bool
		expected    TestResponse
	}{
		{
			name:        "Successful request",
			url:         server.URL + "/success",
			body:        TestRequest{Data: "test data"},
			headers:     nil,
			expectError: false,
			expected: TestResponse{
				Message: "Success",
				Status:  "OK",
				Echo:    "test data",
			},
		},
		{
			name:        "Request with custom headers",
			url:         server.URL + "/custom-header",
			body:        TestRequest{Data: "test data with header"},
			headers:     map[string]string{"X-Custom-Header": "custom-value"},
			expectError: false,
			expected: TestResponse{
				Message: "Custom header received",
				Status:  "OK",
				Echo:    "test data with header",
			},
		},
		{
			name:        "Server error",
			url:         server.URL + "/error",
			body:        TestRequest{Data: "error data"},
			headers:     nil,
			expectError: false, // The function doesn't return an error for non-2xx status codes
			expected: TestResponse{
				Message: "Error",
				Status:  "Error",
			},
		},
		{
			name:        "Invalid URL",
			url:         "http://invalid-url-that-does-not-exist.example",
			body:        TestRequest{Data: "invalid url"},
			headers:     nil,
			expectError: true,
			expected:    TestResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response TestResponse
			err := utils.HttpPostJSON(test.url, test.body, &response, test.headers)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if response.Message != test.expected.Message || response.Status != test.expected.Status {
					t.Errorf("Expected response %+v, got %+v", test.expected, response)
				}

				if test.expected.Echo != "" && response.Echo != test.expected.Echo {
					t.Errorf("Expected echo %s, got %s", test.expected.Echo, response.Echo)
				}
			}
		})
	}
}

func TestHttpPutJSON(t *testing.T) {
	type TestRequest struct {
		Data string `json:"data"`
	}

	type TestResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Echo    string `json:"echo,omitempty"`
	}

	// Setup mock server
	server := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}

		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Check custom header if present
		if r.URL.Path == "/custom-header" && r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header to be custom-value, got %s", r.Header.Get("X-Custom-Header"))
		}

		// Parse request body
		var request TestRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return response based on path
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Success",
				Status:  "OK",
				Echo:    request.Data,
			})
		case "/custom-header":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Custom header received",
				Status:  "OK",
				Echo:    request.Data,
			})
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Error",
				Status:  "Error",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	tests := []struct {
		name        string
		url         string
		body        TestRequest
		headers     map[string]string
		expectError bool
		expected    TestResponse
	}{
		{
			name:        "Successful request",
			url:         server.URL + "/success",
			body:        TestRequest{Data: "test data"},
			headers:     nil,
			expectError: false,
			expected: TestResponse{
				Message: "Success",
				Status:  "OK",
				Echo:    "test data",
			},
		},
		{
			name:        "Request with custom headers",
			url:         server.URL + "/custom-header",
			body:        TestRequest{Data: "test data with header"},
			headers:     map[string]string{"X-Custom-Header": "custom-value"},
			expectError: false,
			expected: TestResponse{
				Message: "Custom header received",
				Status:  "OK",
				Echo:    "test data with header",
			},
		},
		{
			name:        "Server error",
			url:         server.URL + "/error",
			body:        TestRequest{Data: "error data"},
			headers:     nil,
			expectError: false, // The function doesn't return an error for non-2xx status codes
			expected: TestResponse{
				Message: "Error",
				Status:  "Error",
			},
		},
		{
			name:        "Invalid URL",
			url:         "http://invalid-url-that-does-not-exist.example",
			body:        TestRequest{Data: "invalid url"},
			headers:     nil,
			expectError: true,
			expected:    TestResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response TestResponse
			err := utils.HttpPutJSON(test.url, test.body, &response, test.headers)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if response.Message != test.expected.Message || response.Status != test.expected.Status {
					t.Errorf("Expected response %+v, got %+v", test.expected, response)
				}

				if test.expected.Echo != "" && response.Echo != test.expected.Echo {
					t.Errorf("Expected echo %s, got %s", test.expected.Echo, response.Echo)
				}
			}
		})
	}
}

func TestHttpDeleteJSON(t *testing.T) {
	type TestResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	// Setup mock server
	server := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}

		// Check headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type header to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Check custom header if present
		if r.URL.Path == "/custom-header" && r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header to be custom-value, got %s", r.Header.Get("X-Custom-Header"))
		}

		// Return response based on path
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Success",
				Status:  "OK",
			})
		case "/custom-header":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Custom header received",
				Status:  "OK",
			})
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TestResponse{
				Message: "Error",
				Status:  "Error",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	tests := []struct {
		name        string
		url         string
		headers     map[string]string
		expectError bool
		expected    TestResponse
	}{
		{
			name:        "Successful request",
			url:         server.URL + "/success",
			headers:     nil,
			expectError: false,
			expected: TestResponse{
				Message: "Success",
				Status:  "OK",
			},
		},
		{
			name:        "Request with custom headers",
			url:         server.URL + "/custom-header",
			headers:     map[string]string{"X-Custom-Header": "custom-value"},
			expectError: false,
			expected: TestResponse{
				Message: "Custom header received",
				Status:  "OK",
			},
		},
		{
			name:        "Server error",
			url:         server.URL + "/error",
			headers:     nil,
			expectError: false, // The function doesn't return an error for non-2xx status codes
			expected: TestResponse{
				Message: "Error",
				Status:  "Error",
			},
		},
		{
			name:        "Invalid URL",
			url:         "http://invalid-url-that-does-not-exist.example",
			headers:     nil,
			expectError: true,
			expected:    TestResponse{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response TestResponse
			err := utils.HttpDeleteJSON(test.url, &response, test.headers)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if response.Message != test.expected.Message || response.Status != test.expected.Status {
					t.Errorf("Expected response %+v, got %+v", test.expected, response)
				}
			}
		})
	}
}

func TestHttpDownloadFile(t *testing.T) {
	// Setup mock server
	server := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Return response based on path
		switch r.URL.Path {
		case "/download":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("file content"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	tests := []struct {
		name        string
		url         string
		expectError bool
		expected    []byte
	}{
		{
			name:        "Successful download",
			url:         server.URL + "/download",
			expectError: false,
			expected:    []byte("file content"),
		},
		{
			name:        "Server error",
			url:         server.URL + "/error",
			expectError: false, // The function doesn't return an error for non-2xx status codes
			expected:    []byte{},
		},
		{
			name:        "Invalid URL",
			url:         "http://invalid-url-that-does-not-exist.example",
			expectError: true,
			expected:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content, err := utils.HttpDownloadFile(test.url)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				if !reflect.DeepEqual(content, test.expected) {
					t.Errorf("Expected content %s, got %s", test.expected, content)
				}
			}
		})
	}
}

func TestHttpUploadFile(t *testing.T) {
	// Create a temporary file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test-file.txt")
	err := os.WriteFile(tempFile, []byte("test file content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Setup mock server
	mockServer := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check that the request is multipart
		contentType := r.Header.Get("Content-Type")
		if contentType == "" || contentType[:9] != "multipart" {
			t.Errorf("Expected multipart Content-Type, got %s", contentType)
		}

		// Check custom header if present
		if r.URL.Path == "/custom-header" && r.Header.Get("X-Custom-Header") != "custom-value" {
			t.Errorf("Expected X-Custom-Header to be custom-value, got %s", r.Header.Get("X-Custom-Header"))
		}

		// Parse multipart form
		err := r.ParseMultipartForm(10 << 20) // 10 MB
		if err != nil {
			t.Errorf("Error parsing multipart form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check file field
		_, fileHeader, err := r.FormFile("file")
		if err != nil {
			t.Errorf("Error getting file from form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Check filename
		if fileHeader.Filename != "test-file.txt" {
			t.Errorf("Expected filename test-file.txt, got %s", fileHeader.Filename)
		}

		// Check additional fields
		if r.FormValue("field1") != "value1" {
			t.Errorf("Expected field1=value1, got %s", r.FormValue("field1"))
		}

		// Return response based on path
		switch r.URL.Path {
		case "/success", "/custom-header":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("File uploaded successfully"))
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	tests := []struct {
		name             string
		url              string
		fieldName        string
		filePath         string
		additionalFields map[string]string
		headers          map[string]string
		expectError      bool
		expectedStatus   int
	}{
		{
			name:             "Successful upload",
			url:              mockServer.URL + "/success",
			fieldName:        "file",
			filePath:         tempFile,
			additionalFields: map[string]string{"field1": "value1"},
			headers:          nil,
			expectError:      false,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Upload with custom headers",
			url:              mockServer.URL + "/custom-header",
			fieldName:        "file",
			filePath:         tempFile,
			additionalFields: map[string]string{"field1": "value1"},
			headers:          map[string]string{"X-Custom-Header": "custom-value"},
			expectError:      false,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "Server error",
			url:              mockServer.URL + "/error",
			fieldName:        "file",
			filePath:         tempFile,
			additionalFields: map[string]string{"field1": "value1"},
			headers:          nil,
			expectError:      false,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name:             "Invalid URL",
			url:              "http://invalid-url-that-does-not-exist.example",
			fieldName:        "file",
			filePath:         tempFile,
			additionalFields: map[string]string{"field1": "value1"},
			headers:          nil,
			expectError:      true,
			expectedStatus:   0,
		},
		{
			name:             "Non-existent file",
			url:              mockServer.URL + "/success",
			fieldName:        "file",
			filePath:         filepath.Join(tempDir, "non-existent-file.txt"),
			additionalFields: map[string]string{"field1": "value1"},
			headers:          nil,
			expectError:      true,
			expectedStatus:   0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := utils.HttpUploadFile(test.url, test.fieldName, test.filePath, test.additionalFields, test.headers)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else {
					if resp.StatusCode != test.expectedStatus {
						t.Errorf("Expected status code %d, got %d", test.expectedStatus, resp.StatusCode)
					}
					resp.Body.Close()
				}
			}
		})
	}
}
