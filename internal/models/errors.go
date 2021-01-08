package models

import (
	"fmt"
)

// ErrorMessage - answer with error
type ErrorMessage struct {
	Message string `json:"error"`
}

var (
	// ErrInternal -
	ErrInternal = fmt.Errorf("internal error")
)
