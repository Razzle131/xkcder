package core

import "errors"

var (
	ErrBadArguments      = errors.New("arguments are not acceptable")
	ErrAlreadyExists     = errors.New("resource or task already exists")
	ErrNotFound          = errors.New("resource is not found")
	ErrBadResponseStatus = errors.New("bad status of response")
	ErrAlreadyUpdating   = errors.New("update is going on already")
	ErrDbAdapter         = errors.New("db adapter error")
)
