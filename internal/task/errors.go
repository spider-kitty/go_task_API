package task

import "errors"

var (
	ErrTaskNotFound       = errors.New("task not found")
	ErrInvalidTaskID      = errors.New("invalid task id")
	ErrInvalidRequestBody = errors.New("invalid request body")

	ErrTitleRequired       = errors.New("title is required")
	ErrTitleTooShort       = errors.New("title must be at least 3 characters")
	ErrTitleTooLong        = errors.New("title must be less than 100 characters")
	ErrDescriptionTooLong  = errors.New("description must be less than 500 characters")
	ErrCategoryTooLong     = errors.New("category must be less than 50 characters")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidSearchFilter = errors.New("search filter must be less than 100 characters")
)
