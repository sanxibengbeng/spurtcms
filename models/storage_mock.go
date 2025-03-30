package models

import "errors"

// Mock errors
var (
	ErrNotFound     = errors.New("not found")
	ErrStorageError = errors.New("storage error")
)
