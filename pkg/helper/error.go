package helper

import (
	"fmt"
)

type HTTPError struct {
	Code    int
	Message error 
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message.Error())
}

func NewHTTPError(code int, message error) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}
