package validator

import "regexp"

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// checks whether the map has a content or not
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// adds an error message to the map
func (v *Validator) AddError(key string, message string) {
	_, exists := v.Errors[key]
	if !exists {
		v.Errors[key] = message
	}
}

// If the check is failed, adds an error
func (v *Validator) Check(ok bool, key string, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// returns true if the value is exist in a list of string
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// returns true if a value match the regex
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.Match([]byte(value))
}

// unique returns true if all string in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, val := range values {
		uniqueValues[val] = true
	}

	return len(uniqueValues) == len(values)
}
