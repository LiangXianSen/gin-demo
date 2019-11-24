package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// Error is an internal error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	stack   error
}

// Error implements error interface
func (e *Error) Error() string {
	return e.Message
}

// Error returns true if error belongs BadRequest
func (e *Error) IsBadRequest() bool {
	for _, code := range BadRequest {
		if e.Code == code {
			return true
		}
	}
	return false
}

// Wrap returns true if error belongs BadRequest
func (e *Error) Wrap(err error) *Error {
	stack := errors.WithStack(err)
	message := fmt.Sprintf(e.Message + " : " + err.Error())
	return NewErrorWithStack(e.Code, message, stack)
}

func (e *Error) Stack() error {
	return e.stack
}

// NewError returns Error with runtime stack
func NewErrorWithStack(code int, message string, stack error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		stack:   stack,
	}
}

// NewError returns Error instance
func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error declarations
var (
	InvalidParamError     = NewError(InvalidParamErrorCode, "parameter is invalid")
	PermissionDeniedError = NewError(PermissionDeniedErrorCode, "permission denied")
	InvalidTokenError     = NewError(InvalidTokenErrorCode, "invalid token")
	TokenExpiredError     = NewError(TokenExpiredErrorCode, "token expired")
	SignError             = NewError(SignErrorCode, "sign error")
	InternalError         = NewError(InternalErrorCode, "internal error")
)

// Error Code declarations
const (
	InvalidParamErrorCode     = 1002
	PermissionDeniedErrorCode = 1003
	InvalidTokenErrorCode     = 1004
	TokenExpiredErrorCode     = 1005
	SignErrorCode             = 1006
	InternalErrorCode         = 1007
)

// BadRequest is a built-in bad request error list
var BadRequest = []int{
	InvalidParamErrorCode,
	PermissionDeniedErrorCode,
	InvalidTokenErrorCode,
	TokenExpiredErrorCode,
	SignErrorCode,
}
