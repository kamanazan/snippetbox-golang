package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

// standar email regex defined by w3c https://html.spec.whatwg.org/multipage/input.html#valid-e-mail-address
var emailRx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(key, errMsg string) {
	if v.FieldErrors == nil {
		v.FieldErrors = map[string]string{}
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = errMsg
	}
}

func (v *Validator) CheckField(ok bool, key, errMsg string) {
	if !ok {
		v.AddFieldError(key, errMsg)
	}
}

func (v *Validator) AddNonFieldError(errMsg string) {
	v.NonFieldErrors = append(v.NonFieldErrors, errMsg)
}

func StringNotEmpty(val string) bool {
	return strings.TrimSpace(val) != ""
}

func StringInLimit(val string, limit int) bool {
	return utf8.RuneCountInString(val) <= limit
}

func ValueInRange(val int, val2 []int) bool {
	for _, a := range val2 {
		if val == a {
			return true
		}
	}
	return false
}

// MinChars() returns true if a value contains at least n characters.
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Matches() returns true if a value matches a provided compiled regular
// expression pattern.
// func Matches(value string, rx *regexp.Regexp) bool {
// 	return rx.MatchString(value)
// }

func ValidEmail(email string) bool {
	return emailRx.MatchString(email)
}
