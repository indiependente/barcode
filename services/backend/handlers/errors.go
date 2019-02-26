package handlers

import "github.com/pkg/errors"

var (
	ErrMissingCode = errors.New("Missing code parameter")
)
