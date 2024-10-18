package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors    map[string]string
	NonFieldErrors []string
}

var EmailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (v *Validator) IsValid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(value string) {
	v.NonFieldErrors = append(v.NonFieldErrors, value)
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func ValidPassword(value string) bool {
	if !MinChars(value, 8) {
		return false
	}

	validPasswordMap := make(map[string]bool)
	validPasswordMap["hasNumbers"] = false
	validPasswordMap["hasLower"] = false
	validPasswordMap["hasUpper"] = false

	numbers := [10]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	for _, num := range numbers {
		if strings.Contains(value, num) {
			validPasswordMap["hasNumbers"] = true
			break
		}
	}

	lowercaseLetters := [26]string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}

	for _, letter := range lowercaseLetters {
		if strings.Contains(value, letter) {
			validPasswordMap["hasLower"] = true
			break
		}
	}

	uppercaseLetters := [26]string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	for _, letter := range uppercaseLetters {
		if strings.Contains(value, letter) {
			validPasswordMap["hasUpper"] = true
			break
		}
	}

	for _, v := range validPasswordMap {
		if v == true {
			continue
		} else {
			return false
		}
	}
	return true
}
