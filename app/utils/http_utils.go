package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HttpBuildURL constructs a URL from a base URL and query parameters
func HttpBuildURL(baseURL string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// HttpGetJSON performs a GET request and unmarshals the JSON response into the provided interface
func HttpGetJSON(url string, target interface{}, headers map[string]string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// HttpPostJSON performs a POST request with JSON body and unmarshals the response into the provided interface
func HttpPostJSON(url string, body interface{}, target interface{}, headers map[string]string) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// HttpDownloadFile downloads a file from the specified URL
func HttpDownloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// HttpIsSuccessStatusCode checks if the HTTP status code is in the 2xx range
func HttpIsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// HttpPutJSON performs a PUT request with JSON body and unmarshals the response into the provided interface
func HttpPutJSON(url string, body interface{}, target interface{}, headers map[string]string) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// HttpDeleteJSON performs a DELETE request and unmarshals the response into the provided interface
func HttpDeleteJSON(url string, target interface{}, headers map[string]string) error {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

// HttpUploadFile uploads a file to the specified URL using multipart/form-data
func HttpUploadFile(url, fieldName, filePath string, additionalFields map[string]string, headers map[string]string) (*http.Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Add additional form fields
	for key, value := range additionalFields {
		err = writer.WriteField(key, value)
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// Set content type header
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: time.Second * 30, // Longer timeout for file uploads
	}

	return client.Do(req)
}

// HttpCreateHTTPClient creates an HTTP client with custom timeout and transport options
func HttpCreateHTTPClient(timeout time.Duration, maxIdleConns, maxIdleConnsPerHost, maxConnsPerHost int) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		MaxConnsPerHost:     maxConnsPerHost,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

// HttpParseQueryParams parses URL query parameters into a map
func HttpParseQueryParams(queryString string) (map[string]string, error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for key, vals := range values {
		if len(vals) > 0 {
			result[key] = vals[0]
		}
	}

	return result, nil
}
