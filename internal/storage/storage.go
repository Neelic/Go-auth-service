package storage

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrUrlNotFound = errors.New("url not found")
	ErrUrlNotExist = errors.New("url not exist")
)
