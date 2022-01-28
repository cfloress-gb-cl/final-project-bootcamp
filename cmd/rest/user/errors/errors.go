package user

import "errors"

var (
	ErrNotFound          error = errors.New("user not found")
	ErrInternalFailure   error = errors.New("bad request")
	ErrInvalidInput      error = errors.New("invalid fields in request")
	ErrUserAlreadyExists error = errors.New("User already exists")
	ErrNoChangesUpdate   error = errors.New("no changes in data, no records were updated")
	ErrNoMissingField    error = errors.New("Field is not valid or missing")
)

type AppError struct {
	errorMessage error
}

func WrapError(msg error) AppError {
	return AppError{
		errorMessage: msg,
	}
}

func (a AppError) error() error {
	return a.errorMessage
}
