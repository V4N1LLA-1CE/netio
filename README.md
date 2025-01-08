# netio ["net-eye-oh"]
A stable package that simplifies common web server development tasks in Go. Write less boilerplate and focus on building your application logic.

## Overview
I made this package to serve as my own toolkit of reusable components and utilities that I've found myself repeatedly implementing in Go webserver projects. Instead of rewriting these common patterns for each new project, this package provides a centralised, maintained collection of helpers.

As new common patterns emerge in my web development work, I'll continue to add and refine utilities in this package. Feel free to suggest additions or improvements.

## Installation
```bash
go get github.com/V4N1LLA-1CE/netio@latest
```

## Notes:
- **`Netio.Write()`** automatically sets the following headers by default; Make sure to override these using your own headers if needed.

```go
w.Header().Set("Content-Type", "application/json")
w.Header().Set("X-Content-Type-Options", "nosniff")
w.Header().Set("X-Frame-Options", "DENY")
```

## Usage
#### Read JSON from Request 
```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Username string         `json:"username"`
        UserData map[string]any `json:"user_data"`
    }

    // read data into input struct
    err := netio.Read(w, r, &input)
    if err != nil {
        // handle error
    }
}
```
#### Write JSON to Response
```go
func exampleHandler(w http.ResponseWriter, r *http.Request) {
    responseData := map[string]any{
        "message": "Success",
        "data": map[string]interface{}{
            "id": 123,
            "name": "John Doe",
            "email": "john@example.com",
        },
    }
    // headers can be set to  nil for netio.Write()
    // if no further configuration is needed
    headers := http.Header{
        "X-Some-Header": []string{"Value-1", "Value-2"},
        "X-API-Version": []string{"1.0"},
    }
    headers.Add("X-New-Header", "New-Header-Value")

    // netio.Write() will automatically set json headers (can be overriden by custom header)
    err = netio.Write(w, http.StatusOK, netio.Envelope{"example response": responseData}, headers)
    if err != nil {
        // handle error
    }
}
```
#### Validators
```go
v := netio.NewValidator()

email := "invalid-email"
age := 15
role := "superuser"
interests := []string{"coding", "coding"}

emailRx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

v.Check(netio.Matches(email, emailRx), "email", "Invalid email")
v.Check(age >= 18, "age", "Must be 18 or older")
v.Check(netio.IsIn(role, "admin", "user", "moderator"), "role", "Invalid role")
v.Check(!netio.HasDuplicates(interests), "interests", "Duplicate interests found")

if !v.Valid() {
   fmt.Println("Errors:", v.Errors)
}
```

#### JSON HTTP Errors
```go
// read JSON input
var input struct {
    Email string `json:"email"`
    Age   int    `json:"age"`
}

// read request body into input struct
if err := netio.Read(w, r, &input); err != nil {
    netio.Error(w, "error", err, http.StatusBadRequest, nil)
    return
}

// validate input
v := netio.NewValidator()
emailRx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
v.Check(netio.Matches(input.Email, emailRx), "email", "invalid email format")
v.Check(input.Age >= 18, "age", "must be over 18")

if !v.Valid() {
    netio.Error(w, "error", nil, http.StatusUnprocessableEntity, v)
    return
}
```
#### Error response examples
```bash
# Request with no validation errors
curl -X POST 'localhost:8080/register' \
  -H "Content-Type: application/json" \
  -d '{invalid json'

# Response (400 Bad Request)
{
    "error": {
        "status": 400,
        "status_text": "Bad Request",
        "message": "netio.Read(): invalid character 'i' looking for beginning of object key string",
        "timestamp": "2025-01-08T18:45:33.536576+11:00"
    }
}
```

```bash
# Request with validation errors
curl -X POST 'localhost:8080/register' \
  -H "Content-Type: application/json" \
  -d '{
    "email": "notanemail",
    "age": 16
  }'

# Response (422 Unprocessable Entity)
{
    "error": {
        "status": 422,
        "status_text": "Unprocessable Entity",
        "message": "validation failed",
        "validation": {
            "email": "invalid email format",
            "age": "must be over 18"
        },
        "timestamp": "2025-01-08T18:46:33.536576+11:00"
    }
}
```

