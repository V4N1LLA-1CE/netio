package netio

import (
	"regexp"
	"testing"
)

func TestNewValidator(t *testing.T) {
	v := NewValidator()
	if v == nil {
		t.Error("NewValidator() returned nil")
	}
	if v.Errors == nil {
		t.Error("NewValidator() returned validator with nil Errors map")
	}
	if len(v.Errors) != 0 {
		t.Error("NewValidator() returned validator with non-empty Errors map")
	}
}

func TestValidator_AddError(t *testing.T) {
	v := NewValidator()

	// test adding first error
	v.AddError("field1", "error1")
	if len(v.Errors) != 1 {
		t.Error("AddError() failed to add error")
	}
	if v.Errors["field1"] != "error1" {
		t.Error("AddError() stored incorrect error message")
	}

	// test adding duplicate error (should not override)
	v.AddError("field1", "error2")
	if v.Errors["field1"] != "error1" {
		t.Error("AddError() incorrectly overrode existing error")
	}
}

func TestValidator_Valid(t *testing.T) {
	v := NewValidator()

	// Test initial state
	if !v.Valid() {
		t.Error("Valid() returned false for new validator")
	}

	// Test with error
	v.AddError("field1", "error1")
	if v.Valid() {
		t.Error("Valid() returned true for validator with errors")
	}
}

func TestValidator_Check(t *testing.T) {
	tests := []struct {
		name      string
		condition bool
		key       string
		message   string
		wantError bool
	}{
		{"passing condition", true, "field1", "error1", false},
		{"failing condition", false, "field2", "error2", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v := NewValidator()
			v.Check(tc.condition, tc.key, tc.message)

			_, hasError := v.Errors[tc.key]
			if hasError != tc.wantError {
				t.Errorf("Check() error = %v, wantError %v", hasError, tc.wantError)
			}
		})
	}
}

func TestHasDuplicates(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected bool
	}{
		{"no duplicates", []string{"a", "b", "c"}, false},
		{"has duplicates", []string{"a", "b", "a"}, true},
		{"empty slice", []string{}, false},
		{"single value", []string{"a"}, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := HasDuplicates(tc.values); got != tc.expected {
				t.Errorf("HasDuplicates() = %v, want %v", got, tc.expected)
			}
		})
	}

	// test with different types
	numbers := []int{1, 2, 2, 3}
	if !HasDuplicates(numbers) {
		t.Error("HasDuplicates() failed to detect duplicate integers")
	}
}

func TestIsIn(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		allowed  []string
		expected bool
	}{
		{"value present", "a", []string{"a", "b", "c"}, true},
		{"value absent", "d", []string{"a", "b", "c"}, false},
		{"empty allowed values", "a", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIn(tt.value, tt.allowed...); got != tt.expected {
				t.Errorf("IsIn() = %v, want %v", got, tt.expected)
			}
		})
	}

	// test with different types
	if !IsIn(1, 1, 2, 3) {
		t.Error("IsIn() failed with integer type")
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		pattern  string
		expected bool
	}{
		{"matches pattern", "test123", `^[a-z]+\d+$`, true},
		{"no match", "test", `^\d+$`, false},
		{"empty string", "", `^$`, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rx := regexp.MustCompile(tc.pattern)
			if got := Matches(tc.value, rx); got != tc.expected {
				t.Errorf("Matches() = %v, want %v", got, tc.expected)
			}
		})
	}
}
