package myerr

import "errors"

var (
	ErrEmpty   = errors.New("empty data")
	ErrStop    = errors.New("service stopped")
	ErrNotWork = errors.New("service not working")
)
