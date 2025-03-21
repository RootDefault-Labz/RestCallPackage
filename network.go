package network

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var httpClient = &http.Client{
	Transport: &http.Transport{
		// Configure your tls.Config here (RootCAs or InsecureSkipVerify)
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Set to true to bypass verification (not recommended)
			//RootCAs:
		},
	},
}


// logRequest logs the details of the request with a timestamp.
func logRequest(method, endpoint, description string, headers map[string]string, payload string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Print(DottedSeparator)
	log.Printf(LogFormat, timestamp, LogRequestDesc, description)
	log.Printf(LogFormat, timestamp, LogHttpMethod, method)
	log.Printf(LogFormat, timestamp, LogDestEndpoint, endpoint)
	if payload != "" {
		log.Printf(LogFormat, timestamp, LogPayload, payload)
	} else {
		log.Printf(LogFormat, timestamp, LogPayload, LogNullValue)
	}
	log.Printf(LogFormat, timestamp, LogHeaders, "")
	for key, value := range headers {
		log.Printf(LogFormat, timestamp, key, value)
	}
	log.Print(DottedSeparator)
}

// logResponse logs the details of the response with a timestamp.
func logResponse(description string, response string, statusCode int) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Print(DottedSeparator)
	log.Printf(LogFormat, timestamp, LogResponseDesc, description)
	if statusCode != 0 {
		log.Printf(LogFormatInt, timestamp, LogResponseStatus, statusCode)
	}
	if response != "" {
		log.Printf(LogFormat, timestamp, LogResponse, response)
	} else {
		log.Printf(LogFormat, timestamp, LogResponse, LogNullValue)
	}
	log.Print(DottedSeparator)
}

// Add a common request handler
func makeRequest(method, description, urlStr string, payload map[string]interface{}, headers map[string]string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Methods that typically don't have a request body should use query parameters
	isQueryParamMethod := method == methodGET || method == methodDELETE || method == methodHEAD || method == methodOPTIONS

	if isQueryParamMethod && payload != nil {
		q := u.Query()
		for key, value := range payload {
			q.Set(key, fmt.Sprint(value))
		}
		u.RawQuery = q.Encode()
	}

	// Prepare request body for methods that typically have one
	var body io.Reader
	var payloadStr string
	if !isQueryParamMethod && payload != nil {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}
		body = bytes.NewBuffer(jsonPayload)
		payloadStr = string(jsonPayload)
	}

	// Create the request
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return "", err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Log the request details
	logRequest(method, u.String(), description, headers, payloadStr)

	// Perform the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := ReadResponseBody(resp)
	if err != nil {
		return "", err
	}

	// Log the response details
	logResponse(description, responseBody, resp.StatusCode)

	// Check for non-2xx status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("received non-2xx response code: %d, response:%v", resp.StatusCode, responseBody)
	}

	return responseBody, nil
}

// Update the public functions to use the common handler
func MakeGETRequest(description, baseURL string, queryParams map[string]string, headers map[string]string) (string, error) {
	payload := make(map[string]interface{})
	for k, v := range queryParams {
		payload[k] = v
	}
	return makeRequest(methodGET, description, baseURL, payload, headers)
}

func MakePOSTRequest(description, url string, payload map[string]interface{}, headers map[string]string) (string, error) {
	return makeRequest(methodPOST, description, url, payload, headers)
}

func MakePUTRequest(description, url string, payload map[string]interface{}, headers map[string]string) (string, error) {
	return makeRequest(methodPUT, description, url, payload, headers)
}

func MakeDELETERequest(description, url string, queryParams map[string]string, headers map[string]string) (string, error) {
	payload := make(map[string]interface{})
	for k, v := range queryParams {
		payload[k] = v
	}
	return makeRequest(methodDELETE, description, url, payload, headers)
}

func MakePATCHRequest(description, url string, payload map[string]interface{}, headers map[string]string) (string, error) {
	return makeRequest(methodPATCH, description, url, payload, headers)
}

func MakeHEADRequest(description, url string, queryParams map[string]string, headers map[string]string) (string, error) {
	payload := make(map[string]interface{})
	for k, v := range queryParams {
		payload[k] = v
	}
	return makeRequest(methodHEAD, description, url, payload, headers)
}

func MakeOPTIONSRequest(description, url string, queryParams map[string]string, headers map[string]string) (string, error) {
	payload := make(map[string]interface{})
	for k, v := range queryParams {
		payload[k] = v
	}
	return makeRequest(methodOPTIONS, description, url, payload, headers)
}

// ReadResponseBody simplified to remove duplicate defer
func ReadResponseBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
