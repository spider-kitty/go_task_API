package task

import "strings"

func validateTitle(title string) error {
	if title == "" {
		return ErrTitleRequired
	}

	if len(title) < 3 {
		return ErrTitleTooShort
	}

	if len(title) > 100 {
		return ErrTitleTooLong
	}

	return nil
}

func validateDescription(description string) error {
	if len(description) > 500 {
		return ErrDescriptionTooLong
	}

	return nil
}

func validateCategory(category string) error {
	if len(category) > 50 {
		return ErrCategoryTooLong
	}

	return nil
}

func normalizeStatus(status string) (string, error) {
	status = strings.ToLower(strings.TrimSpace(status))

	if status == "" {
		return "todo", nil
	}

	if status != "todo" && status != "in_progress" && status != "done" {
		return "", ErrInvalidStatus
	}

	return status, nil
}

func validateSearchFilter(search string) error {
	if len(search) > 100 {
		return ErrInvalidSearchFilter
	}

	return nil
}
