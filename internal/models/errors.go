package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrDuplicateUsername  = errors.New("models: duplicate username")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrInvalidCredentails = errors.New("models: invalid creadentails")
)
