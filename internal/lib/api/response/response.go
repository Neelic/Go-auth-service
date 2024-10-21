package response

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func Ok() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errsMessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errsMessages = append(errsMessages, err.Param()+" is required")
		case "url":
			errsMessages = append(errsMessages, "field "+err.Param()+" is required")
		default:
			errsMessages = append(errsMessages, err.Param()+" is invalid")
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errsMessages, ", "),
	}
}
