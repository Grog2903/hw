package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInvalidRule     = errors.New("invalid rule")
	ErrInvalidParam    = errors.New("invalid validation param")
	ErrBigValue        = errors.New("value bigger than max")
	ErrSmallValue      = errors.New("value less than max")
	ErrNotInValue      = errors.New("value not available")
	ErrIncorrectLength = errors.New("incorrect length")
	ErrRegexp          = errors.New("invalid regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var result strings.Builder
	for _, err := range v {
		result.WriteString(fmt.Sprintf("%s: %s\n", err.Field, err.Err.Error()))
	}

	return result.String()
}

func (v ValidationErrors) Unwrap() error {
	if len(v) == 0 {
		return nil
	}
	return v[0].Err
}

func (v ValidationErrors) Is(target error) bool {
	for _, validationError := range v {
		if errors.Is(validationError.Err, target) {
			return true
		}
	}

	return false
}

func Validate(v interface{}) (ValidationErrors, error) {
	var ValidationErrors ValidationErrors
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return nil, nil
	}

	st := reflect.TypeOf(v)
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		validateTag := field.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		value := val.Field(i)

		rules := strings.Split(validateTag, "|")
		for _, rule := range rules {
			validationErr, err := applyRule(rule, value)
			if validationErr != nil {
				ValidationErrors = append(ValidationErrors, ValidationError{Field: field.Name, Err: validationErr})
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return ValidationErrors, nil
}

func applyRule(rule string, value reflect.Value) (error, error) {
	ruleParts := strings.SplitN(rule, ":", 2)
	if len(ruleParts) != 2 {
		return ErrInvalidRule, nil
	}
	validateMethod, validateParam := ruleParts[0], ruleParts[1]

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.String:
		return validateString(value.String(), validateMethod, validateParam)
	case reflect.Int:
		return validateInt(int(value.Int()), validateMethod, validateParam)
	case reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			validationErr, err := applyRule(rule, value.Index(i))
			if err != nil {
				return nil, err
			}
			if validationErr != nil {
				return validationErr, nil
			}
		}
	default:
		return nil, nil
	}

	return nil, nil
}

func validateInt(value int, method string, param string) (error, error) {
	switch method {
	case "max":
		maxValue, err := strconv.Atoi(param)
		if err != nil {
			return nil, err
		}
		if value > maxValue {
			return ErrBigValue, nil
		}
	case "min":
		minValue, err := strconv.Atoi(param)
		if err != nil {
			return nil, err
		}
		if value < minValue {
			return ErrSmallValue, nil
		}
	case "in":
		options := strings.Split(param, ",")
		for _, option := range options {
			optionValue, err := strconv.Atoi(option)
			if err != nil {
				return nil, err
			}
			if optionValue == value {
				return nil, nil
			}
		}

		return ErrNotInValue, nil
	}

	return nil, nil
}

func validateString(value string, method string, param string) (error, error) {
	switch method {
	case "len":
		length, err := strconv.Atoi(param)
		if err != nil {
			return nil, err
		}

		if len(value) != length {
			return ErrIncorrectLength, nil
		}
	case "regexp":
		re, err := regexp.Compile(param)
		if err != nil {
			return nil, err
		}
		if !re.MatchString(value) {
			return ErrRegexp, nil
		}
	case "in":
		options := strings.Split(param, ",")
		for _, option := range options {
			if value == option {
				return nil, nil
			}
		}
		return ErrNotInValue, nil
	}

	return nil, nil
}
