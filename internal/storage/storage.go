package storage

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlExist    = errors.New("url not exist")
)
