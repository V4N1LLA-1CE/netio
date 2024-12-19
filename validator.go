package netio

import (
	"regexp"
	"slices"
)

// Validator provides a structure for collecting and managing validation errors.
// It maintains a map of field-specific error messages that can be accumulated
// during the validation process.
type Validator struct {
	Errors map[string]string
}

// NewValidator is a helper function that creates and initializes a new
// Validator instance with an empty error map.
// This is the recommended way to create a new Validator.
//
// Example:
//
//	v := netio.NewValidator()
//	v.Check(user.Age >= 18, "age", "must be at least 18 years old")
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// Valid returns true if the validator has no errors, false otherwise.
// This method can be used to determine if all validation checks have passed.
//
// Example:
//
//	v := netio.NewValidator()
//	v.Check(user.Age >= 18, "age", "must be at least 18 years old")
//
//	if !v.Valid() {
//	    // Handle validation errors
//	}
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a specific field to the validator's error map.
// If an error already exists for the given key, it will not be overwritten.
func (v *Validator) AddError(key, message string) {
	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = message
	}
}

// Check performs a validation check based on a condition. If the condition is false,
// it adds an error message for the specified key.
//
// Parameters:
//   - condition: The boolean condition to check
//   - key: The field or identifier for the potential error
//   - message: The error message to store if the condition is false
//
// Example:
//
//	v := netio.NewValidator()
//	v.Check(len(password) >= 8, "password", "must be at least 8 characters")
func (v *Validator) Check(condition bool, key, message string) {
	if !condition {
		v.AddError(key, message)
	}
}

// HasDuplicates checks if a slice contains any duplicate values.
// It uses Go's generics to work with any comparable type.
//
// Parameters:
//   - values: A slice of comparable values to check for duplicates
//
// Returns:
//   - bool: true if duplicates are found, false otherwise
//
// Example:
//
//	v := netio.NewValidator()
//	ids := []int{1, 2, 2, 3}
//	v.Check(!netio.HasDuplicates(ids), "id", "Duplicate IDs found")
//	if !v.Valid() {
//	  // handle validation err
//	}
func HasDuplicates[T comparable](values []T) bool {
	uniqueVals := make(map[T]bool)

	for _, val := range values {
		if _, exists := uniqueVals[val]; exists {
			return true
		}

		uniqueVals[val] = true
	}

	return false
}

// IsIn checks if a value exists within a set of allowed values.
// It uses Go's generics to work with any comparable type.
//
// Parameters:
//   - value: The value to check
//   - allowedValues: A variadic parameter of allowed values to check against
//
// Returns:
//   - bool: true if the value is found in allowedValues, false otherwise
//
// Example:
//
//	role := "administrator"
//	v := netio.NewValidator()
//	v.Check(netio.IsIn(role, "admin", "user", "moderator"), "role", "Invalid role")
//	if !v.Valid() {
//	  // handle validation err
//	}
func IsIn[T comparable](value T, allowedValues ...T) bool {
	return slices.Contains(allowedValues, value)
}

// Matches checks if a string value matches a given regular expression pattern.
// It uses Go's regexp package for pattern matching.
//
// Parameters:
//   - value: The string to check
//   - rx: A compiled regular expression pattern
//
// Returns:
//   - bool: true if the string matches the pattern, false otherwise
//
// Example:
//
//	emailRx := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
//	v := netio.NewValidator()
//	v.Check(netio.Matches(email, emailRx), "email", "Invalid email format")
//	if !v.Valid() {
//	    // handle validation err
//	}
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
