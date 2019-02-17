package handlers

import "github.com/pkg/errors"

var (
	ErrMissingCode error = errors.New("Missing code parameter")
)
