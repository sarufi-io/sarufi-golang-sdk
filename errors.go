package sarufi

import (
	"fmt"
	"strings"
)

type (
	RequestError struct {
		Detail  string `json:"detail,omitempty"`
		Message string `json:"message,omitempty"`
	}

	ValidationErrorDetail struct {
		Loc  []string `json:"loc"`
		Msg  string   `json:"msg"`
		Type string   `json:"type"`
	}

	ValidationError struct {
		Detail []ValidationErrorDetail `json:"detail"`
	}
)

func (re *RequestError) Error() string {
	if re != nil {
		return fmt.Sprintf("request error: %s", strings.Join([]string{re.Detail, re.Message}, ", "))
	}
	return ""
}

// String returns the string representation of the ValidationErrorDetail
func (v *ValidationErrorDetail) String() string {
	locations := strings.Join(v.Loc, ",")
	return fmt.Sprintf("location: [%s], message: %s, error type: %s", locations, v.Msg, v.Type)
}

// FormatValidationError formats the validation error into a string
func (e *ValidationError) Error() string {
	stringErrs := make([]string, 0)
	for _, v := range e.Detail {
		stringErrs = append(stringErrs, v.String())
	}
	return fmt.Sprintf("validation error: %s", strings.Join(stringErrs, ", "))
}
