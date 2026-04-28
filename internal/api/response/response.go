package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	// Status codes
	StatusOK    = "ok"
	StatusError = "error"

	// Validation error tags
	ValidationErrorUrl      = "url"
	ValidationErrorRequired = "required"
)

func OK() Response {
	return Response{Status: StatusOK}
}

func Error(msg string) Response {
	return Response{Status: StatusError, Error: msg}
}

func ValidationError(errs validator.ValidationErrors) Response {
	// Convert the validation errors to a slice of strings
	var errMsgs []string
	// Loop over the validation errors and append the error messages to the slice
	for _, err := range errs {
		switch err.ActualTag() {
		case ValidationErrorUrl:
			errMsgs = append(errMsgs, fmt.Sprintf("invalid URL: %s", err.Param()))
		case ValidationErrorRequired:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("unknown validation error: %s", err.ActualTag()))
		}
	}
	//return Response with the error messages
	return Response{Status: StatusError, Error: strings.Join(errMsgs, ", ")}
}
