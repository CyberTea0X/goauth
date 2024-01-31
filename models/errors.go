package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type ConstantError string

func (e ConstantError) Error() string { return string(e) }

// constant errors because i don't like sentinel errors
const (
	ErrNoTokenSpecified = ConstantError("no token specified")
	ErrInvalidToken     = ConstantError("invalid token")
	ErrTokenExpired     = ConstantError("token expired")
	ErrInvalidJson      = ConstantError("invalid JSON")
)

// I want to return the same status code and message as the external service
type ExternalServiceError struct {
	Status int
	Msg    string
}

func (e *ExternalServiceError) Error() string {
	if e.Msg == "" {
		return "external service error"
	} else {
		return e.Msg
	}
}

func (e *ExternalServiceError) Is(err error) bool {
	_, ok := err.(*ExternalServiceError)
	return ok
}

func NewExternalServiceError(status int) *ExternalServiceError {
	return &ExternalServiceError{status, ""}
}

type errorConverter interface {
	Error() string
}

func ErrToMap(err errorConverter) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

type errorJson struct {
	Error string `json:"error"`
}

// parse error from json {"error":"error"} or error string
// nil if error is not parsed
// Function does not close response body
func ErrFromResponse(response *http.Response) error {
	if response.ContentLength == -1 {
		return nil
	}
	content := response.Header.Get("Content-Type")
	switch content {
	case "application/json":
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil
		}
		e := new(errorJson)
		if err := json.Unmarshal(body, e); err != nil {
			return nil
		}
		return errors.New(e.Error)
	case "text/plain":
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil
		}
		return errors.New(string(body))
	default:
		return nil
	}
}
