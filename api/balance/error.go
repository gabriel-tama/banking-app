package balance

import "errors"

var (
	ErrValidationFailed = errors.New("validation req failed")
	ErrInvalidToken     = errors.New("token invalid or missing")
	ErrNotEnoughBalance = errors.New("not enough balances")
)
