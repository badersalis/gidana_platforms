package services

import "errors"

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string { return e.Message }

func ErrBadRequest(msg string) error    { return &ServiceError{Code: 400, Message: msg} }
func ErrUnauthorized(msg string) error  { return &ServiceError{Code: 401, Message: msg} }
func ErrForbidden(msg string) error     { return &ServiceError{Code: 403, Message: msg} }
func ErrNotFound(msg string) error      { return &ServiceError{Code: 404, Message: msg} }
func ErrInternal(msg string) error      { return &ServiceError{Code: 500, Message: msg} }

func IsServiceError(err error) (*ServiceError, bool) {
	var se *ServiceError
	if errors.As(err, &se) {
		return se, true
	}
	return nil, false
}
