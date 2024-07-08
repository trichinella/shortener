package repo

import (
	"fmt"
	"shortener/internal/app/entity"
)

type DuplicateShortcutError struct {
	Shortcut *entity.Shortcut
	Err      error
}

func (e DuplicateShortcutError) Error() string {
	return fmt.Errorf("shortcut %s already has been created with ID: %s. %w", e.Shortcut.OriginalURL, e.Shortcut.ID, e.Err).Error()
}

func NewDuplicateShortcutError(err error, shortcut *entity.Shortcut) error {
	return &DuplicateShortcutError{
		Shortcut: shortcut,
		Err:      err,
	}
}
