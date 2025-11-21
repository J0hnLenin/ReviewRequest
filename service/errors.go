package service

import "errors"

var (
	ErrConnection      = errors.New("repository connection error")
	ErrQueryExecution  = errors.New("error during query execution")
)