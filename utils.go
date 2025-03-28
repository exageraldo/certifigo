package certifigo

import (
	"errors"
	"reflect"
	"regexp"
)

// Bool is a type alias for a pointer to a boolean value (*bool).
// This type is used to differentiate between a boolean value that is explicitly set
// (true or false) and one that is left undefined (nil). By using a pointer, we can
// determine whether a value was intentionally provided or omitted.
//
// This distinction is particularly important when working with functions like [Merge],
// where default values can cause unintended behavior. For example, in Go, the default
// value of a boolean variable is `false`. If we do not handle this explicitly using
// pointers, a new `false` value will never be recognized as explicitly set, because
// it is indistinguishable from the default value.
//
// To facilitate the use of Bool, it is common to define True and False variables
// as pointers to their respective boolean values. These variables can then be
// used to assign values to Bool without needing to create new pointers each time.
//
// Example:
//
//	var myBool Bool
//	myBool = True  // Assigns true
//	myBool = False // Assigns false
//	myBool = nil   // Represents an undefined state
type Bool *bool

// In Go, you cannot directly take the address of a literal value like `false`.
// To work around this limitation, an array literal is used to hold the value,
// and the address of the first element of the array is taken.
var True = &[]bool{true}[0]
var False = &[]bool{false}[0]

var (
	ErrNoMatchFound       = errors.New("the value provided does not match the expected expression")
	ErrInvalidConfigRegex = errors.New("the regular expression provided is invalid")
	ErrGroupMustBeUnique  = errors.New("regex group names must be unique")
)

// FindNamedMatches matches the given inline configuration string with the provided regular expression pattern.
// It returns a map of key-value pairs where the keys are the named subexpressions in the pattern and the values are the corresponding matches.
// If no match is found, it returns an error.
//
// Reference:
//   - https://stackoverflow.com/a/20751656
func FindNamedMatches(str, exp string) (map[string]string, error) {
	configRegex, err := regexp.Compile(exp)
	if err != nil {
		return nil, ErrInvalidConfigRegex
	}

	match := configRegex.FindStringSubmatch(str)
	// when the match is empty, it means no match was found
	if len(match) == 0 {
		return nil, ErrNoMatchFound
	}

	// grouping the values found to their keys
	result := make(map[string]string)
	for i, name := range configRegex.SubexpNames() {
		if i != 0 && name != "" {
			if _, ok := result[name]; ok {
				return nil, ErrGroupMustBeUnique
			}
			result[name] = match[i]
		}
	}

	// even if the match was found, the result can still be empty
	// if the regular expression does not have named subexpressions
	if len(result) == 0 {
		return nil, ErrNoMatchFound
	}

	return result, nil
}

// Merge takes two values of the same type and merges the fields of the second value (override)
// into the first value (base). Only non-zero fields from the override value are copied to the base value.
// This function uses reflection to iterate over the fields of the struct, so it works with any struct type.
//
// Type Parameters:
//   - T: The type of the values to be merged. It can be any type.
//
// Parameters:
//   - base: The base value that will be updated with non-zero fields from the override value.
//   - override: The value whose non-zero fields will be used to update the base value.
//
// Returns:
//   - T: A new value of type T with the merged fields.
//
// Note:
//   - This function skips unexported fields (fields that cannot be set).
//   - It uses reflection, which may have performance implications and should be used with caution.
func Merge[T any](base, override T) T {
	baseVal := reflect.ValueOf(&base).Elem()
	overrideVal := reflect.ValueOf(override)

	for i := range baseVal.NumField() {
		field := baseVal.Field(i)
		overrideField := overrideVal.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Check if the field in override has a non-zero value
		zero := reflect.Zero(overrideField.Type())
		isZeroValue := reflect.DeepEqual(overrideField.Interface(), zero.Interface())
		if !isZeroValue {
			field.Set(overrideField)
		}
	}

	return base
}
