package netio

import (
	"regexp"
	"slices"
)

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exist := v.Errors[key]; !exist {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(condition bool, key, message string) {
	if !condition {
		v.AddError(key, message)
	}
}

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

func IsIn[T comparable](value T, allowedValues ...T) bool {
	return slices.Contains(allowedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
