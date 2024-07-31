package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var result strings.Builder
	var repeatableValue string
	var isSlashes bool

	for i, currentRune := range input {
		err := checkError(input, i, repeatableValue)
		if err != nil {
			return "", err
		}

		if repeatableValue == "\\" && !isSlashes {
			result.WriteString(string(currentRune))
			repeatableValue = string(currentRune)
			isSlashes = true
			continue
		}

		if unicode.IsDigit(currentRune) {
			repeatableCountInt, _ := strconv.Atoi(string(currentRune))

			if repeatableCountInt == 0 {
				str := result.String()
				result.Reset()
				result.WriteString(str[:len(str)-1])
				repeatableValue = ""
			} else {
				result.WriteString(strings.Repeat(repeatableValue, repeatableCountInt-1))
				repeatableValue = ""
			}

			continue
		}

		if string(currentRune) != "\\" {
			result.WriteString(string(currentRune))
		}

		repeatableValue = string(currentRune)
		isSlashes = false
	}

	return result.String(), nil
}

func checkError(input string, incr int, repeatableValue string) error {
	currentRune := rune(input[incr])
	isCurrentDigit := unicode.IsDigit(currentRune)
	isNotLastChar := incr != len(input)-1
	isNextDigit := isNotLastChar && unicode.IsDigit(rune(input[incr+1]))
	isPreviousCharValid := incr-1 > 0 && input[incr-1] != '\\'

	if incr == 0 && isCurrentDigit {
		return ErrInvalidString
	}

	if isCurrentDigit && isNextDigit && isPreviousCharValid {
		return ErrInvalidString
	}

	if repeatableValue == "\\" && !(isCurrentDigit || input[incr] == '\\') {
		return ErrInvalidString
	}

	return nil
}
