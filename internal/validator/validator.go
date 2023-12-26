package validator

import (
	"strings"
	"unicode/utf8"
)


type Validator struct {
    FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
    return len(v.FieldErrors) == 0
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