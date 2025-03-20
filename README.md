# RestCallPackage

A robust Go package that simplifies making HTTP REST API calls with built-in logging and error handling.

## Features

- Support for all standard HTTP methods:
  - GET
  - POST
  - PUT
  - DELETE
  - PATCH
  - HEAD
  - OPTIONS
- Automatic request/response logging
- Query parameter handling
- JSON payload support
- Custom header support
- TLS configuration
- Comprehensive error handling

## Usage

### Making a GET Request

```go
response, err := MakeGETRequest(
    "Get Users",                              // Description
    "https://api.example.com/users",          // URL
    map[string]string{"page": "1"},           // Query Parameters
    map[string]string{"Accept": "application/json"}, // Headers
)
```

### Making a POST Request

```go
payload := map[string]interface{}{
    "name": "John Doe",
    "email": "john@example.com",
}

response, err := MakePOSTRequest(
    "Create User",                            // Description
    "https://api.example.com/users",          // URL
    payload,                                  // Request Body
    map[string]string{"Content-Type": "application/json"}, // Headers
)
```

### Other HTTP Methods

The package supports similar syntax for other HTTP methods:

- `MakePUTRequest()`
- `MakeDELETERequest()`
- `MakePATCHRequest()`
- `MakeHEADRequest()`
- `MakeOPTIONSRequest()`

## Logging

The package automatically logs:
- Request details (method, URL, headers, payload)
- Response details (status code, response body)
- Timestamps for all operations

## Error Handling

- Automatic handling of non-2xx status codes
- URL parsing validation
- JSON marshalling/unmarshalling
- Response body reading

## TLS Configuration

The package includes configurable TLS settings for secure communications.

## License

[Add your license information here]
