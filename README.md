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
### Read JSON from Request
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
### Write JSON to Response
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
### Validators for Data Validation
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

### JSON Http Error Wrapper
```go
err := netio.Read(w, r, &input)
if err != nil {
    netio.Error(w, err.Error(), http.StatusUnprocessableEntity)
    return
}
```
