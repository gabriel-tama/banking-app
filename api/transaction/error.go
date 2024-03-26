package transaction

import "errors"

var (
	ErrValidationFailed = errors.New("validation req failed")
	ErrInvalidToken     = errors.New("token invalid or missing")
)
