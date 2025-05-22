package repositories

import "errors"

var (
	ErrorNotFound = errors.New("Not Found")
	ErrorInvalid  = errors.New("Invalid")
)
