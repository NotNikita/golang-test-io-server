package model

import "errors"

// Tasks Repo possible errors
var (
	ErrStorageNil        = errors.New("storage is nil")
	ErrTaskMapNil        = errors.New("task map is nil")
	ErrInvalidTask       = errors.New("invalid input task")
	ErrTaskAlreadyExists = errors.New("task with this ID already exists")
	ErrTaskNotFound      = errors.New("task not found")
)
