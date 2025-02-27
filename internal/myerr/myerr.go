package myerr

import "errors"

var (
	ErrEmpty = errors.New("Empty data")
	ErrStop  = errors.New("Service stopped")
)
