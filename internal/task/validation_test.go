package task

import (
	"errors"
	"testing"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantError error
	}{
		{
			name:      "valid title",
			title:     "Learn Go",
			wantError: nil,
		},
		{
			name:      "empty title",
			title:     "",
			wantError: ErrTitleRequired,
		},
		{
			name:      "too short title",
			title:     "Go",
			wantError: ErrTitleTooShort,
		},
		{
			name:      "too long title",
			title:     string(make([]byte, 101)),
			wantError: ErrTitleTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTitle(tt.title)

			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantError   error
	}{
		{
			name:        "valid description",
			description: "This is a valid description",
			wantError:   nil,
		},
		{
			name:        "empty description is valid",
			description: "",
			wantError:   nil,
		},
		{
			name:        "too long description",
			description: string(make([]byte, 501)),
			wantError:   ErrDescriptionTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.description)

			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateCategory(t *testing.T) {
	tests := []struct {
		name      string
		category  string
		wantError error
	}{
		{
			name:      "valid category",
			category:  "backend",
			wantError: nil,
		},
		{
			name:      "empty category is valid",
			category:  "",
			wantError: nil,
		},
		{
			name:      "too long category",
			category:  string(make([]byte, 51)),
			wantError: ErrCategoryTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCategory(tt.category)

			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		wantStatus string
		wantError  error
	}{
		{
			name:       "empty status returns default todo",
			status:     "",
			wantStatus: "todo",
			wantError:  nil,
		},
		{
			name:       "valid todo status",
			status:     "todo",
			wantStatus: "todo",
			wantError:  nil,
		},
		{
			name:       "valid done status with spaces and uppercase",
			status:     " DONE ",
			wantStatus: "done",
			wantError:  nil,
		},
		{
			name:       "valid in progress status",
			status:     "in_progress",
			wantStatus: "in_progress",
			wantError:  nil,
		},
		{
			name:       "invalid status",
			status:     "finished",
			wantStatus: "",
			wantError:  ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := normalizeStatus(tt.status)

			if status != tt.wantStatus {
				t.Errorf("expected status %q, got %q", tt.wantStatus, status)
			}

			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}

func TestValidateSearchFilter(t *testing.T) {
	tests := []struct {
		name      string
		search    string
		wantError error
	}{
		{
			name:      "valid search",
			search:    "go",
			wantError: nil,
		},
		{
			name:      "empty search is valid",
			search:    "",
			wantError: nil,
		},
		{
			name:      "too long search",
			search:    string(make([]byte, 101)),
			wantError: ErrInvalidSearchFilter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSearchFilter(tt.search)

			if !errors.Is(err, tt.wantError) {
				t.Errorf("expected error %v, got %v", tt.wantError, err)
			}
		})
	}
}
